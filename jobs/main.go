package jobs

import (
	"context"
	"encoding/json"
	"karmaclips/config"
	"karmaclips/database"
	"karmaclips/helpers/generations"
	"karmaclips/segmind"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func StartJobs() {
	for true {
		time.Sleep(5 * time.Second)
		go segmindImageJob()
	}
}

func segmindImageJob() {
	opt, err := redis.ParseURL(config.NewConfig().RedisURL)
	if err != nil {
		log.Printf("Error parsing Redis URL: %v", err)
		return
	}

	client := redis.NewClient(opt)
	ctx := context.Background()

	for {
		keys, err := client.Keys(ctx, "karmaclips:*").Result()
		if err != nil {
			log.Printf("Error getting keys: %v", err)
			return
		}

		pendingJobsExist := false

		for _, key := range keys {
			jobDataStr, err := client.Get(ctx, key).Result()
			if err != nil {
				log.Printf("Error getting job data for key %s: %v", key, err)
				continue
			}

			var jobData JobData
			err = json.Unmarshal([]byte(jobDataStr), &jobData)
			if err != nil {
				log.Printf("Error unmarshaling job data for key %s: %v", key, err)
				continue
			}

			if jobData.Status == "pending" {
				pendingJobsExist = true

				// Update status to "processing" and save back to Redis
				jobData.Status = "processing"
				jobDataBytes, err := json.Marshal(jobData)
				if err != nil {
					log.Printf("Error marshalling job data for key %s: %v", key, err)
					continue
				}

				err = client.Set(ctx, key, jobDataBytes, 0).Err()
				if err != nil {
					log.Printf("Error setting job data for key %s: %v", key, err)
					continue
				}

				// Process the job
				imageUri, err := segmind.RequestCreateImage(jobData.Prompt, jobData.Model, jobData.BatchSize, jobData.Width, jobData.Height)
				if err != nil {
					log.Printf("Error creating image for key %s: %v", key, err)
					// Update status to "error"
					jobData.Status = "error"
					// jobData.Error Msg = err.Error()
					jobDataBytes, _ = json.Marshal(jobData)
					client.Set(ctx, key, jobDataBytes, 0)
					continue
				}

				// Update job data with image URL and status "completed"
				jobData.Url = *imageUri
				jobData.Status = "completed"

				//Push to DB
				generation := &database.Generation{
					CreatedBy:   jobData.UserId,
					CreditsUsed: jobData.CreditsUsed,
					Timestamp:   time.Now(),
					MediaUri:    jobData.Url,
					Type:        jobData.Type,
					Meta: database.Meta{
						ModelId:        jobData.Model,
						Dimensions:     string(jobData.Height) + "x" + string(jobData.Width),
						Prompt:         jobData.Prompt,
						NegativePrompt: jobData.NegativePrompt,
						BatchSize:      jobData.BatchSize,
					},
				}

				_, err = generations.CreateGeneration(generation)
				if err != nil {
					log.Fatal("DB failure")
				}

				// jobDataBytes, err = json.Marshal(jobData)
				// if err != nil {
				// 	log.Printf("Error marshalling job data for key %s: %v", key, err)
				// 	continue
				// }

				// err = client.Set(ctx, key, jobDataBytes, 0).Err()
				// if err != nil {
				// 	log.Printf("Error setting job data for key %s: %v", key, err)
				// 	continue
				// }
				client.Del(ctx, key)
			}
		}

		if !pendingJobsExist {
			// No pending jobs found, exit the loop
			break
		}
	}
}

type JobData struct {
	Prompt         string `json:"prompt"`
	BatchSize      int    `json:"batch_size"`
	Model          string `json:"model"`
	Status         string `json:"status"`
	Url            string `json:"url"`
	Height         int    `json:"height"`
	Width          int    `json:"width"`
	UserId         string `json:"user_id"`
	CreditsUsed    int    `json:"credits_used"`
	Type           string `json:"type"`
	NegativePrompt string `json:"negative_prompt"`
}

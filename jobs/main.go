package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"karmaclips/config"
	"karmaclips/segmind"
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
		return
	}

	client := redis.NewClient(opt)

	ctx := context.Background()
	jobList := client.Keys(ctx, "karmaclips:*")
	jobListResult, err := jobList.Result()
	if len(jobListResult) < 1 {
		return
	}
	i := 0
	do := jobListResult[i]
	var jd JobData
	for i < len(jobListResult) {
		i += 1
		jobData, err := client.Get(ctx, do).Result()
		if err != nil {
			return
		}
		err = json.Unmarshal([]byte(jobData), &jd)
		if err != nil {
			return
		}
		if jd.Status == "pending" {
			break
		}
	}
	jd.Status = "processing"

	jdBytes, err := json.Marshal(jd)
	if err != nil {
		fmt.Println("Error marshalling jd:", err)
		return
	}
	client.Set(ctx, do, jdBytes, 0)

	reqImageUri, err := segmind.RequestCreateImage(jd.Prompt, jd.Model, jd.BatchSize, jd.Height, jd.Width)
	if err != nil {
		fmt.Println("Error creating image")
		return
	}

	jd.Url = *reqImageUri
	jd.Status = "completed"
	jdBytes, err = json.Marshal(jd)
	if err != nil {
		fmt.Println("Error marshalling jd:", err)
		return
	}
	client.Set(ctx, do, jdBytes, 0)
	return
}

type JobData struct {
	Prompt    string `json:"prompt"`
	BatchSize int    `json:"batch_size"`
	Model     string `json:"model"`
	Status    string `json:"status"`
	Url       string `json:"url"`
	Height    int    `json:"height"`
	Width     int    `json:"width"`
}

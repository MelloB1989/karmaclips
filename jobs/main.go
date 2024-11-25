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
	fmt.Println("Starting job")
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
	do := jobListResult[0]
	jobData, err := client.Get(ctx, do).Result()
	if err != nil {
		return
	}
	var jd JobData
	err = json.Unmarshal([]byte(jobData), &jd)
	// client.Del(ctx, do)
	jd.Status = "processing"
	client.Set(ctx, do, jd, time.Second*2400)

	reqImageUri, err := segmind.RequestCreateImage(jd.Prompt, jd.Model, jd.BatchSize, jd.Height, jd.Width)
	if err != nil {
		fmt.Println("Error creating image")
		return
	}

	jd.Url = *reqImageUri
	client.Set(ctx, do, jd, time.Second*2400)
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

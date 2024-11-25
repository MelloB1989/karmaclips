package gen

import (
	"context"
	"encoding/json"
	"fmt"
	"karmaclips/config"
	"karmaclips/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type CreateImageReq struct {
	Prompt    string `json:"prompt"`
	BatchSize int    `json:"batch_size"`
	Model     string `json:"model"`
	Height    int    `json:"height"`
	Width     int    `json:"width"`
}

func CreateImage(c *fiber.Ctx) error {
	req := new(CreateImageReq)
	if err := c.BodyParser(req); err != nil {
		fmt.Println("Error parsing request")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	var genJob = map[string]interface{}{
		"prompt":     req.Prompt,
		"batch_size": req.BatchSize,
		"model":      req.Model,
		"height":     req.Height,
		"width":      req.Width,
		"status":     "pending",
		"url":        "",
	}

	jobId := "karmaclips:" + utils.GenerateID()

	opt, _ := redis.ParseURL(config.NewConfig().RedisURL)
	client := redis.NewClient(opt)

	genJobJson, err := json.Marshal(genJob)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error creating job",
			"message": err,
			"data":    nil,
		})
	}
	client.Set(context.Background(), jobId, genJobJson, 0)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Job created",
		"data":    jobId,
	})
}

func GetJobStatus(c *fiber.Ctx) error {
	jobId := c.Params("jobId")
	var genJob map[string]interface{}

	opt, err := redis.ParseURL(config.NewConfig().RedisURL)
	if err != nil {
		log.Println(err)
	}
	client := redis.NewClient(opt)

	ctx := context.Background()
	genJobJson, err := client.Get(ctx, jobId).Result()
	if err != nil {
		log.Println("Error getting order from Redis:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error getting job",
			"message": err,
			"data":    nil,
		})
	}

	err = json.Unmarshal([]byte(genJobJson), &genJob)
	if err != nil {
		log.Println("Error unmarshalling order:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error getting job",
			"message": err,
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Job created",
		"data":    genJob,
	})
}

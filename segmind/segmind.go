package segmind

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"karmaclips/config"
	"net/http"
)

func RequestCreateImage(prompt string, batch_size int, width int, height int) {
	data := map[string]interface{}{
		"prompt":          prompt,
		"negative_prompt": "low quality, blurry",
		"steps":           25,
		"guidance_scale":  5.5,
		"seed":            98552302,
		"sampler":         "euler",
		"scheduler":       "sgm_uniform",
		"width":           width,
		"height":          height,
		"aspect_ratio":    "custom",
		"batch_size":      batch_size,
		"image_format":    "jpeg",
		"image_quality":   95,
		"base64":          true,
	}
	api := config.NewConfig().SegmindSDAPI
	jsonPayload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error converting struct to json:", err)
		return
	}
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.NewConfig().SegmindAPIKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

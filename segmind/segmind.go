package segmind

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"karmaclips/config"
	"net/http"
	"os"
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
	api := config.NewConfig().SegmindSamaritanAPI
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

	// Parse JSON response to get the base64 image data
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	// Get the base64 encoded image string
	imageData, ok := responseData["image"].(string)
	if !ok {
		fmt.Println("No image data found in response")
		return
	}

	// Decode the base64 image data
	imageBytes, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		fmt.Println("Error decoding base64 image data:", err)
		return
	}

	// Save the image to a file
	fileName := "generated_image1.jpeg"
	err = os.WriteFile(fileName, imageBytes, 0644)
	if err != nil {
		fmt.Println("Error saving image file:", err)
		return
	}

	fmt.Println("Image saved successfully as:", fileName)
}

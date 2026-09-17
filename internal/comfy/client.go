package comfy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	url        string
	httpClient *http.Client
}

func New(url string) (*Client, error) {
	err := isReachable(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		url:        url,
		httpClient: &http.Client{},
	}, nil

}

func isReachable(rawURL string) error {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Head(rawURL)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	return nil
}

func (c *Client) SendWorkflow(workflow []byte) (string, error) {
	resp, err := c.httpClient.Post(
		c.url+"/prompt",
		"application/json",
		bytes.NewReader(workflow),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf("Response: %s", body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ComfyUI returned %s: %s", resp.Status, body)
	}

	var result struct {
		PromptID   string         `json:"prompt_id"`
		Number     int            `json:"number"`
		NodeErrors map[string]any `json:"node_errors"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	return result.PromptID, nil
}

func (c *Client) WaitForResultAndGetName(promptID string) (string, error) {
	var emptyBody uint8 = 0
	for {
		time.Sleep(2 * time.Second)
		resp, err := c.httpClient.Get(c.url + "/history/" + promptID)
		if err != nil {
			return "", err
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		resp.Body.Close()
		var result map[string]any

		err = json.Unmarshal(body, &result)
		if err != nil {
			return "", err
		}

		if len(result) == 0 {
			if emptyBody > 30 {
				return "", fmt.Errorf("Empty 30 times")
			}
			emptyBody++
		}

		prompt, ok := result[promptID].(map[string]any)
		if !ok {
			continue
		}

		status, ok := prompt["status"].(map[string]any)
		if !ok {
			continue
		}

		completed, ok := status["completed"].(bool)
		if !ok {
			continue
		}

		if completed {
			log.Println("Generation completed!")
			filename, err := c.findFilename(result, promptID)
			if err != nil {
				return "", err
			}
			return filename, nil
		}
	}
}

func (c *Client) findFilename(jsonBodyMapped map[string]any, promptID string) (string, error) {
	prompt, ok := jsonBodyMapped[promptID].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid prompt structure")
	}

	outputs, ok := prompt["outputs"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid outputs structure")
	}

	saveNode, ok := outputs["10052"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid saveNode structure")
	}

	imagesList, ok := saveNode["images"].([]any)
	if !ok || len(imagesList) == 0 {
		return "", fmt.Errorf("images is not a slice or empty")
	}

	firstImage, ok := imagesList[0].(map[string]any)
	if !ok {
		return "", fmt.Errorf("invalid image item structure")
	}

	filename, ok := firstImage["filename"].(string)
	if !ok {
		return "", fmt.Errorf("filename not found or not a string")
	}

	return filename, nil
}

func (c *Client) GetImageBytes(filename string) ([]byte, error) {
	resp, err := c.httpClient.Get(c.url + "/view?filename=" + filename)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetImageBytes: comfyui returned bad status: %s", resp.Status)
	}

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read comfy response: %w", err)
	}
	return imgBytes, nil
}

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

func New(url string) *Client {
	return &Client{
		url:        url,
		httpClient: &http.Client{},
	}
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

func (c *Client) WaitForResult(promptID string) (string, error) {
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
			if emptyBody > 10 {
				return "", fmt.Errorf("Empty 10 times")
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
			return c.findFilename(promptID), nil
		}
	}
}

func (c *Client) findFilename(result map[string]any) (string, error) {

	return promptID
}

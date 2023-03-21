package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gopherzz/simple-openai-client/internal/models"
)

type OpenAiClient struct {
	ApiToken  string
	Model     string
	MaxTokens int
	Debug     bool
	Logger    *log.Logger
}

// Make request to OpenAi GPT Api
func (c OpenAiClient) MakeOpenAiReq(prompt string) (string, error) {
	endpoint := "https://api.openai.com/v1/completions"

	requestBody, err := json.Marshal(map[string]interface{}{
		"model":      c.Model,
		"prompt":     prompt,
		"max_tokens": c.MaxTokens,
		"n":          1,
	})

	if err != nil {
		log.Println("Error creating request body:", err)
		return "", nil
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}

	if c.Debug {
		c.Logger.Println(req)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ApiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	choices := models.OpenaiResponse{}
	err = json.Unmarshal(body, &choices)
	if err != nil {
		log.Println(err)
		return "", err
	}

	if c.Debug {
		c.Logger.Println(choices)
	}

	return choices.Choices[0].Text, nil
}

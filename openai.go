package openai

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type openaiResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string `json:"text"`
		Index        int    `json:"index"`
		Logprobs     any    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type OpenAiClient struct {
	ApiToken  string
	Model     string
	MaxTokens int
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ApiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	choices := openaiResponse{}
	err = json.Unmarshal(body, &choices)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return choices.Choices[0].Text, nil
}

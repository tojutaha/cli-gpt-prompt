package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
    "encoding/json"
)

type Response struct {
    ID       string `json:"id"`
    Object   string `json:"object"`
    Created  int64  `json:"created"`
    Model    string `json:"model"`
    Choices  []struct {
        Index        int `json:"index"`
        Message      struct {
            Role    string `json:"role"`
            Content string `json:"content"`
        } `json:"message"`
        FinishReason string `json:"finish_reason"`
    } `json:"choices"`
    Usage struct {
        PromptTokens     int `json:"prompt_tokens"`
        CompletionTokens int `json:"completion_tokens"`
        TotalTokens      int `json:"total_tokens"`
    } `json:"usage"`
}

func parseJSONResponse(jsonData string) (string, error) {
    var data Response
    err := json.Unmarshal([]byte(jsonData), &data)
    if err != nil {
        return "", err
    }

    // Assuming theres only one choice in the choices array
    if len(data.Choices) > 0 {
        return data.Choices[0].Message.Content, nil
    }

    return "", fmt.Errorf("No choices found in the response")
}

func queryChatGPT(apiKey string, prompt string) string {
    url := "https://api.openai.com/v1/chat/completions"
    personality := "You are a funny guy. You really like to make puns and dad jokes every now and then and start every response with heh"
    //payload := strings.NewReader(fmt.Sprintf(`{ "model": "gpt-3.5-turbo", "messages": [{"role": "user", "content": "%s"}], "temperature": 0.7}`, prompt))
    payload := strings.NewReader(fmt.Sprintf(`{ "model": "gpt-3.5-turbo", "messages": [{"role": "system", "content": "%s"}, {"role": "user", "content": "%s"}], "temperature": 0.7}`, personality, prompt))
    req, _ := http.NewRequest("POST", url, payload)
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Authorization", "Bearer "+apiKey)

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal("Error making the API request:", err)
    }
    defer res.Body.Close()

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal("Error reading API response", err)
    }

    return string(body)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Printf("Usage: %s \"Your prompt here\"", os.Args[0])
        os.Exit(1)
    }

    // TODO: Better way of handling this..?
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        fmt.Println("The environment variable is not set. Shutting down")
        os.Exit(1)
    }

    prompt := os.Args[1]
    jsonResponse := queryChatGPT(apiKey, prompt)

    text, err := parseJSONResponse(jsonResponse)
    if err != nil {
        fmt.Println("Error parsing JSON response", err)
        return
    }

    fmt.Println(text)
}


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

type Choice struct {
    Text         string `json:"text"`
    Index           int `json:"index"`
    Logprobs     string `json:"logprobs"`
    FinishReason string `json:"finish_reason"`
}

type Response struct {
    Choices []Choice `json:"choices"`
}

func parseJSONResponse(jsonData string) (string, error) {
    var data Response
    err := json.Unmarshal([]byte(jsonData), &data)
    if err != nil {
        return "", err
    }

    // Assuming theres only one choice in the choices array
    if len(data.Choices) > 0 {
        return data.Choices[0].Text, nil
    }

    return "", fmt.Errorf("No choices found in the response")
}

func queryChatGPT(apiKey string, prompt string) string {
    url := "https://api.openai.com/v1/engines/text-davinci-003/completions"
    payload := strings.NewReader(fmt.Sprintf(`{"prompt": "%s", "temperature": 0.5, "max_tokens": 1000}`, prompt))

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
    if len(os.Args[1]) < 2 {
        fmt.Println("Usage: go run main.go \"Your prompt here\"")
        os.Exit(1)
    }

    // TODO: Better way of handling this..
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


package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var model = "gpt-3.5-turbo"
//var personality = "You are a funny guy. You really like to make puns and dad jokes every now and then and start every response like you were being bored of my questions"
var personality = "You are a nature documentary narrator and start every response like you were being bored of my questions"

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

func queryChatGPT(apiKey string, usePersonality bool, prompt string) string {
    url := "https://api.openai.com/v1/chat/completions"
    var payload = strings.NewReader("")
    if usePersonality {
        payload = strings.NewReader(fmt.Sprintf(`{ "model": "%s", "messages": [{"role": "system", "content": "%s"}, {"role": "user", "content": "%s"}], "temperature": 0.7}`, model, personality, prompt))
    } else {
        payload = strings.NewReader(fmt.Sprintf(`{ "model": "%s", "messages": [{"role": "user", "content": "%s"}], "temperature": 0.7}`, model, prompt))
    }

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

    // If 3rd argument exists and is greater than 0 use custom system message
    usePersonality := false
    if len(os.Args) == 3 {
        argStr := os.Args[2]
        arg, err := strconv.Atoi(argStr)
        if err != nil {
            fmt.Println("Error:", err)
            os.Exit(1)
        }

        if arg < 0 {
            fmt.Println("Argument 3 must be a non-negative number")
            os.Exit(1)
        }

        if arg > 0 {
            usePersonality = true
        }
    }

    // TODO: Better way of handling this..?
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        fmt.Println("The environment variable is not set. Shutting down")
        os.Exit(1)
    }

    prompt := os.Args[1]
    jsonResponse := queryChatGPT(apiKey, usePersonality, prompt)

    text, err := parseJSONResponse(jsonResponse)
    if err != nil {
        fmt.Println("Error parsing JSON response", err)
        return
    }

    fmt.Println(text)
}


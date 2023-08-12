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

var messages = []string{}

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

func saveHistory(arr []string) {
    // Convert array to single string with each element separated by a new line
    data := strings.Join(arr, "\n")

    err := ioutil.WriteFile("message_history", []byte(data), 0644)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
}

func loadHistory() []string {
    data, err := ioutil.ReadFile("message_history")
    if err != nil {
        if !os.IsNotExist(err) {
            fmt.Println("Error:", err)
        }
        return nil
    }

    stringArray := strings.Split(string(data), "\n")
    return stringArray
}

func addToHistory(arr* []string, data string) {
    // If theres more than one element in the array, add comma at the end of the string, 
    // so that its valid JSON format
    if len(*arr) > 0 {
        for i := 0; i < len(*arr); i++ {
            // only if there were no comma already
            if !strings.HasSuffix((*arr)[i], ",") {
                (*arr)[i] = (*arr)[i] + ","
            }
        }
    }

    *arr = append(*arr, data)
}

func parseJSONResponse(jsonData string) (string, error) {
    var data Response
    err := json.Unmarshal([]byte(jsonData), &data)
    if err != nil {
        return "", err
    }

    // Assuming theres only one choice in the choices array
    if len(data.Choices) > 0 {
        content := data.Choices[0].Message.Content
        strings.ReplaceAll(content, "\"", "")
        message := fmt.Sprintf(`{"role": "assistant", "content": "%s"}`, content)
        addToHistory(&messages, message)
        return content, nil
    }

    return "", fmt.Errorf("No choices found in the response")
}

func queryChatGPT(apiKey string, prompt string) string {
    url := "https://api.openai.com/v1/chat/completions"

    message := fmt.Sprintf(`{"role": "user", "content": "%s"}`, prompt)
    addToHistory(&messages, message)

    payload := strings.NewReader(fmt.Sprintf(`{ "model": "gpt-3.5-turbo", "messages": %s, "temperature": 0.7}`, messages))
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
    var useSystemMessage = 0
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
            useSystemMessage = 1
        }
    }

    messages = loadHistory()
    if messages == nil {
        if useSystemMessage > 0 {
            //personality = "You are a funny guy. You really like to make puns and dad jokes every now and then and start every response like you were being bored of my questions"
            personality := "You are a nature documentary narrator and start every response like you are bored of my questions"
            message := fmt.Sprintf(`{"role": "system", "content": "%s"}`, personality)
            addToHistory(&messages, message)
        }
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
        fmt.Println(jsonResponse)
        fmt.Println(messages)
        fmt.Println("Error parsing JSON response", err)
        return
    }

    saveHistory(messages)

    fmt.Println(text)
}


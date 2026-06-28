package main

import (
	"bytes"
    "encoding/json"
    "fmt"
    "net/http"
	"os"
)

type ChatRequest struct {
    Message string `json:"message"`
}
type GroqRequest struct {
	Model    string        `json:"model"`
	Messages []GroqMessage `json:"messages"`
}
type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type GroqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
    http.HandleFunc("/chat", handleChat)
    fmt.Println("server running on :8080")
    http.ListenAndServe(":8080", nil)
}

func handleChat(w http.ResponseWriter, r *http.Request) {
    var req ChatRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

	groqReq := GroqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []GroqMessage{
			{Role: "system", Content: "You are a concise technical assistant."},
			{Role: "user", Content: req.Message},
		},
	}

	body, _ := json.Marshal(groqReq)
	httpReq, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		http.Error(w, "groq call failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// step 5: decode groq's response
	var groqResp GroqResponse
	json.NewDecoder(resp.Body).Decode(&groqResp)

	// step 6: return the answer
	fmt.Fprintln(w, groqResp.Choices[0].Message.Content)
    
}

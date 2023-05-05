package model

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type OpenAIRequest struct {
	Model       string  `json:"model"`
	Temperature float32 `json:"temperature"`
	Prompt      string  `json:"prompt"`
}

type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"text-davinci-003"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Text         string   `json:"text"`
	Index        int      `json:"index"`
	Logprobs     Logprobs `json:"logprobs"`
	FinishReason string   `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Logprobs struct {
	LikelyTokens int
	Valid        bool
}

var nullLiteral = []byte("null")

func (l *Logprobs) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, nullLiteral) {
		l.Valid = false
		return nil
	}

	err := json.Unmarshal(b, &l.LikelyTokens)
	if err == nil {
		l.Valid = true
		return nil
	}

	return err
}

func (l Logprobs) MarshalJSON() ([]byte, error) {
	if l.Valid {
		return json.Marshal(l.LikelyTokens)
	} else {
		return nullLiteral, nil
	}
}

func (l Logprobs) String() string {
	if !l.Valid {
		return "null"
	}
	return fmt.Sprintf("%v", l.LikelyTokens)
}

package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateArticleFromTranscript(transcript string) (string, error) {
	ctx := context.Background()

	apiKey, ok := os.LookupEnv("LLM_KEY")
	if !ok {
		log.Fatalln("Environment variable LLM_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash-exp")

	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	session := model.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text("I will send you a transcription of the podcast. Can you generate medium like article from it?\n"+transcript))
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}

	var article strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		article.WriteString(fmt.Sprintf("%v", part))
	}

	return article.String(), nil
}

type QuizQuestion struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   int      `json:"true_answer_index"`
}

type Quiz struct {
	Questions []QuizQuestion `json:"questions"`
}

func GenerateQuizzesFromArticle(article string) (*Quiz, error) {
	ctx := context.Background()

	apiKey, ok := os.LookupEnv("LLM_KEY")
	if !ok {
		log.Fatalln("Environment variable LLM_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	model := client.GenerativeModel("gemini-2.0-flash-exp")

	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"
	session := model.StartChat()
	session.History = []*genai.Content{}

	prompt := fmt.Sprintf(`Generate a 5-question quiz based on this article. Follow these rules:
	1. Each question must have exactly 4 options
	2. Options must be plausible distractors
	3. True answer index must be 0-3
	4. Output must be valid JSON matching this format:
	{
		"questions": [
			{
				"question": "...",
				"options": ["a", "b", "c", "d"],
				"true_answer_index": 0
			}
		]
	}
	Article: %s`, article)

	resp, err := session.SendMessage(ctx, genai.Text(prompt))
	if err != nil {
		return nil, fmt.Errorf("error sending message: %v", err)
	}
	var rawResponse strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		rawResponse.WriteString(fmt.Sprintf("%v", part))
	}

	// Clean up the response by removing markdown code block markers
	cleanedResponse := strings.TrimPrefix(rawResponse.String(), "```json\n")
	cleanedResponse = strings.TrimSuffix(cleanedResponse, "```\n")
	cleanedResponse = strings.TrimSpace(cleanedResponse) // Remove any remaining whitespace

	var quizResp Quiz
	if err := json.Unmarshal([]byte(cleanedResponse), &quizResp); err != nil {
		fmt.Println(cleanedResponse)
		return nil, fmt.Errorf("failed to parse quiz response: %v", err)
	}

	return &quizResp, nil
}

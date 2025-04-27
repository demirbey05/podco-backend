package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateArticleFromTranscript(transcript, language string) (string, error) {
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

	model := client.GenerativeModel("gemini-2.5-flash-preview-04-17")

	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	session := model.StartChat()
	session.History = []*genai.Content{}

	resp, err := session.SendMessage(ctx, genai.Text(generateArticlePrompt(transcript, language)))
	if err != nil {
		return "", fmt.Errorf("error sending message: %v", err)
	}

	var article strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		article.WriteString(fmt.Sprintf("%v", part))
	}
	if strings.Contains(article.String(), `"error"`) {
		// Parse the JSON to extract the error message
		var errorResponse struct {
			Error string `json:"error"`
		}
		err := json.Unmarshal([]byte(article.String()), &errorResponse)
		if err == nil && errorResponse.Error != "" {
			return "", errors.New(errorResponse.Error)
		}
		// Fallback error if parsing fails
		return "", errors.New("the provided content could not be processed as educational material")
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

func GenerateQuizzesFromArticle(article, language string) (*Quiz, error) {
	ctx := context.Background()

	apiKey, ok := os.LookupEnv("LLM_KEY")
	if !ok {
		log.Fatalln("Environment variable LLM_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	model := client.GenerativeModel("gemini-2.5-flash-preview-04-17")

	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"
	session := model.StartChat()
	session.History = []*genai.Content{}

	prompt := generateQuizPrompt(article, language)

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
	cleanedResponse = strings.TrimSuffix(cleanedResponse, "\n")
	cleanedResponse = strings.TrimSuffix(cleanedResponse, "```")
	cleanedResponse = strings.TrimSpace(cleanedResponse) // Remove any remaining whitespace

	var quizResp Quiz
	if err := json.Unmarshal([]byte(cleanedResponse), &quizResp); err != nil {
		fmt.Println(cleanedResponse)
		return nil, fmt.Errorf("failed to parse quiz response: %v", err)
	}

	fmt.Println("Quiz generated successfully")

	return &quizResp, nil
}

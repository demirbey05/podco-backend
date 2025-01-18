package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var model *genai.GenerativeModel

func init() {

	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found")
	}
	ctx := context.Background()

	apiKey, ok := os.LookupEnv("LLM_KEY")
	if !ok {
		log.Fatalln("Environment variable GEMINI_API_KEY not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	model = client.GenerativeModel("gemini-2.0-flash-exp")

	model.SetTemperature(1)
	model.SetTopK(40)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

}

func GenerateArticleFromTranscript(transcript string) (string, error) {
	ctx := context.Background()

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

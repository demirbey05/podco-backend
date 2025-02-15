package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/demirbey05/auth-demo/internal/store"
)

const (
	ArticleGenerated int = iota
	QuizGenerated
	Error
)

func CreateNewPod(link, userID string, podStore store.PodStore) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	// Insert a pod, job and set goroutines

	title, err := GetYouTubeVideoTitle(link)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting video title: %v", err)
	}
	podId, err := podStore.InsertPod(ctx, link, title, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("error inserting pod: %v", err)
	}
	jobId, err := podStore.InsertPodJob(ctx, podId)
	if err != nil {
		return 0, 0, fmt.Errorf("error inserting job: %v", err)
	}

	trans, err := getTranscript(link)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting transcript: %v", err)
	}

	generateArticleJob(trans, podStore, podId, jobId)
	return podId, jobId, nil

}

// TranscriptResponse represents the JSON structure returned by the transcriber service.
type TranscriptResponse struct {
	VideoID    string `json:"video_id"`
	Transcript string `json:"transcript"`
}

func getTranscript(link string) (string, error) {
	// Load TRANSCRIBER_URL from environment variables
	transcriberURL, exists := os.LookupEnv("TRANSCRIBER_URL")
	if !exists {
		return "", fmt.Errorf("TRANSCRIBER_URL is not set in the environment")
	}

	// Construct the full URL with query parameters
	endpoint := fmt.Sprintf("%s/transcript", transcriberURL)
	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("invalid transcriber URL: %v", err)
	}

	// Add the 'url' query parameter
	query := reqURL.Query()
	query.Set("url", link)
	reqURL.RawQuery = query.Encode()

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new HTTP GET request with context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request to transcriber: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("transcriber returned non-OK status: %s", resp.Status)
	}

	// Decode the JSON response
	var transcriptResp TranscriptResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&transcriptResp); err != nil {
		return "", fmt.Errorf("error decoding transcriber response: %v", err)
	}

	return transcriptResp.Transcript, nil
}

func generateArticleJob(transcript string, podStore store.PodStore, podId, jobId int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	article, err := GenerateArticleFromTranscript(transcript)
	if err != nil {
		podStore.UpdatePodJob(ctx, jobId, Error)
		return
	}

	if err := podStore.InsertArticle(ctx, podId, article); err != nil {
		podStore.UpdatePodJob(ctx, jobId, Error)
		return
	}

	if err := podStore.UpdatePodJob(ctx, jobId, ArticleGenerated); err != nil {
		podStore.UpdatePodJob(ctx, jobId, Error)
		return
	}

	generateQuizJob(article, podStore, podId, jobId)
}
func generateQuizJob(article string, podStore store.PodStore, podID, jobId int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	quiz, err := GenerateQuizzesFromArticle(article)
	if err != nil {
		fmt.Println(err)
		podStore.UpdatePodJob(ctx, jobId, Error)
		return
	}

	quizID, err := podStore.InsertQuiz(ctx, podID)
	if err != nil {
		fmt.Println(err)
		podStore.UpdatePodJob(ctx, jobId, Error)
		return
	}

	for _, question := range quiz.Questions {
		if _, err := podStore.InsertQuestion(ctx, quizID, question.Question, question.Options, question.Answer); err != nil {
			fmt.Println(err)
			podStore.UpdatePodJob(ctx, jobId, Error)
			return
		}
	}

	if err := podStore.UpdatePodJob(ctx, jobId, QuizGenerated); err != nil {
		fmt.Println(err)
		podStore.UpdatePodJob(ctx, jobId, Error)
		return
	}

	fmt.Println("Quiz submitted successfully")

}

package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/demirbey05/auth-demo/internal/store"
)

var languageMap = map[string]string{
	"English":    "en",
	"Spanish":    "es",
	"French":     "fr",
	"German":     "de",
	"Italian":    "it",
	"Portuguese": "pt",
	"Dutch":      "nl",
	"Polish":     "pl",
	"Russian":    "ru",
	"Japanese":   "ja",
	"Korean":     "ko",
	"Chinese":    "zh",
	"Turkish":    "tr",
	"Hindi":      "hi",
}

const (
	ArticleGenerated int = iota
	QuizGenerated
	Error
)

func CreateNewPod(link, userID, language string, podStore store.PodStore, usageStore store.UsageStore) (int, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	// Insert a pod, job and set goroutines
	canonLink, err := CanonicalizeYouTubeURL(link)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error canonicalizing link: %v", err)
	}
	link = canonLink
	// Get cost of the job
	duration, err := GetYouTubeVideoDuration(link)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error getting video duration: %v", err)
	}

	cost := int(math.Ceil((duration / 24) * 10))
	fmt.Println("cost is ", cost)
	fmt.Println("duration is ", duration)
	// Check Credit
	remaining, err := usageStore.GetRemainingCredits(ctx, userID)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error getting remaining credits: %v", err)
	}
	if remaining < cost {
		return 0, 0, 0, fmt.Errorf("insufficient credits")
	}

	title, err := GetYouTubeVideoTitle(link)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error getting video title: %v", err)
	}
	podId, err := podStore.InsertPod(ctx, link, title, userID)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error inserting pod: %v", err)
	}
	jobId, err := podStore.InsertPodJob(ctx, podId)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error inserting job: %v", err)
	}
	remaining, err = usageStore.DecrementCredit(ctx, userID, cost)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error decrementing credit: %v", err)
	}

	var trans string
	if os.Getenv("ENV") == "dev" {
		trans, err = getTranscript(link, language)
	} else {
		trans, err = getTranscriptFromAPI(link)
	}
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error getting transcript: %v", err)
	}

	generateArticleJob(trans, language, podStore, podId, jobId)
	return podId, jobId, remaining, nil

}

// TranscriptResponse represents the JSON structure returned by the transcriber service.
type TranscriptResponse struct {
	VideoID    string `json:"video_id"`
	Transcript string `json:"transcript"`
}

func getTranscript(link, language string) (string, error) {
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
	langCode, ok := languageMap[language]
	if !ok {
		fmt.Println(language)
		return "", fmt.Errorf("invalid language")
	}
	// Add the 'url' query parameter
	query := reqURL.Query()
	query.Set("url", link)
	query.Set("language", langCode)
	reqURL.RawQuery = query.Encode()

	// Create an HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Create a new HTTP GET request with context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println(reqURL.String())
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

// TranscriptResponse represents the structure of the API response
type APITranscriptResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

// getTranscript makes an API call to retrieve the transcript and merges the text segments
func getTranscriptFromAPI(videoURL string) (string, error) {
	// Create the request URL
	url := fmt.Sprintf("https://api.supadata.ai/v1/youtube/transcript?url=%s", videoURL)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set the API key in the request header
	req.Header.Set("x-api-key", os.Getenv("SUPADATA_API_KEY"))

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	var transcript APITranscriptResponse
	if err := json.Unmarshal(body, &transcript); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	// Merge the text segments
	var mergedText bytes.Buffer
	for _, segment := range transcript.Content {
		mergedText.WriteString(segment.Text)
		mergedText.WriteString(" ") // Add a space between segments
	}

	return mergedText.String(), nil
}

func generateArticleJob(transcript, language string, podStore store.PodStore, podId, jobId int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	article, err := GenerateArticleFromTranscript(transcript, language)
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

	generateQuizJob(article, language, podStore, podId, jobId)
}
func generateQuizJob(article, language string, podStore store.PodStore, podID, jobId int) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	quiz, err := GenerateQuizzesFromArticle(article, language)
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

// CanonicalizeYouTubeURL converts a YouTube URL (e.g. youtu.be/VIDEO_ID)
// into its canonical form: https://www.youtube.com/watch?v=VIDEO_ID.
func CanonicalizeYouTubeURL(videoURL string) (string, error) {
	u, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf("invalid url: %v", err)
	}

	var videoID string
	switch u.Host {
	case "youtu.be":
		// For URLs like "https://youtu.be/8u2pW2zZLCs?si=0nYPTaKO1--mX0uU"
		videoID = strings.TrimPrefix(u.Path, "/")
	case "www.youtube.com", "youtube.com":
		// For URLs like "https://www.youtube.com/watch?v=8u2pW2zZLCs"
		if u.Path == "/watch" {
			videoID = u.Query().Get("v")
		} else if strings.HasPrefix(u.Path, "/embed/") {
			videoID = strings.TrimPrefix(u.Path, "/embed/")
		}
	default:
		return "", fmt.Errorf("not a youtube url")
	}

	if videoID == "" {
		return "", fmt.Errorf("could not extract video id")
	}

	canonical := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	return canonical, nil
}

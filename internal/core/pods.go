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

func CreateNewPod(link string, podStore store.PodStore) (int, int, error) {
	// Insert a pod, job and set goroutines

	trans, err := getTranscript(link)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting transcript: %v", err)
	}

	fmt.Println(trans)
	return 0, 0, nil

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

package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// YouTubeVideoResponse represents the structure of the YouTube API response
type YouTubeVideoResponse struct {
	Items []struct {
		Snippet struct {
			Title string `json:"title"`
		} `json:"snippet"`
	} `json:"items"`
}

// GetYouTubeVideoTitle fetches the title of a YouTube video using the YouTube Data API
func GetYouTubeVideoTitle(videoURL string) (string, error) {

	apiKey := os.Getenv("YOUTUBE_API_KEY")
	// Extract video ID from the URL
	videoID := strings.Split(videoURL, "v=")[1]
	if strings.Contains(videoID, "&") {
		videoID = strings.Split(videoID, "&")[0]
	}

	// Build the API request URL
	apiURL := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet&id=%s&key=%s", videoID, apiKey)

	// Make the HTTP GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to make API request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse the JSON response
	var videoResponse YouTubeVideoResponse
	if err := json.Unmarshal(body, &videoResponse); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Check if the video was found
	if len(videoResponse.Items) == 0 {
		return "", fmt.Errorf("video not found")
	}

	// Return the video title
	return videoResponse.Items[0].Snippet.Title, nil
}

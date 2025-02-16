package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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

// VideoListResponse is used to parse the YouTube API response.
type VideoListResponse struct {
	Items []struct {
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// parseISO8601Duration parses a subset of ISO8601 durations (e.g. PT3M45S, PT1H2M3S)
// For simplicity, this function supports hours, minutes and seconds.
func parseISO8601Duration(iso string) (time.Duration, error) {
	var hours, minutes, seconds int
	_, err := fmt.Sscanf(iso, "PT%dH%dM%dS", &hours, &minutes, &seconds)
	if err != nil {
		// Try a format without hours
		_, err = fmt.Sscanf(iso, "PT%dM%dS", &minutes, &seconds)
		if err != nil {
			// Try a format with minutes only
			_, err = fmt.Sscanf(iso, "PT%dM", &minutes)
			if err != nil {
				// Try seconds only
				_, err = fmt.Sscanf(iso, "PT%dS", &seconds)
				if err != nil {
					return 0, fmt.Errorf("unsupported duration format: %s", iso)
				}
			}
		}
	}
	dur := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second
	return dur, nil
}

// GetYouTubeVideoDuration fetches the video duration from the YouTube Data API.
func GetYouTubeVideoDuration(canonicalURL string) (float64, error) {
	// Extract the video id from the canonical URL.
	u, err := url.Parse(canonicalURL)
	if err != nil {
		return 0, fmt.Errorf("invalid canonical url: %v", err)
	}
	videoID := u.Query().Get("v")
	if videoID == "" {
		return 0, fmt.Errorf("could not extract video id")
	}

	apiKey, ok := os.LookupEnv("YOUTUBE_API_KEY")
	if !ok {
		return 0, fmt.Errorf("YOUTUBE_API_KEY is not set")
	}

	endpoint := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=contentDetails&id=%s&key=%s", videoID, apiKey)
	resp, err := http.Get(endpoint)
	if err != nil {
		return 0, fmt.Errorf("failed to call YouTube API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("YouTube API returned non-OK status: %s", resp.Status)
	}

	var apiResp VideoListResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, fmt.Errorf("error decoding YouTube API response: %v", err)
	}
	if len(apiResp.Items) == 0 {
		return 0, fmt.Errorf("no video found for id: %s", videoID)
	}

	dur, err := parseISO8601Duration(apiResp.Items[0].ContentDetails.Duration)
	if err != nil {
		return 0, fmt.Errorf("error parsing video duration: %v", err)
	}
	return dur.Minutes(), nil
}

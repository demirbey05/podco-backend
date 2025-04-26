package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

// VideoListResponse is used to parse the YouTube API response.
type VideoListResponse struct {
	Items []struct {
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
	} `json:"items"`
}

// parseISO8601Duration parses ISO8601 duration strings and returns the total duration in minutes.
func parseISO8601Duration(iso string) (int, error) {
	var hours, minutes, seconds int

	// Handle various ISO8601 duration formats
	if strings.Contains(iso, "H") {
		if strings.Contains(iso, "M") && strings.Contains(iso, "S") {
			// Format: PT#H#M#S
			_, err := fmt.Sscanf(iso, "PT%dH%dM%dS", &hours, &minutes, &seconds)
			if err != nil {
				return 0, fmt.Errorf("unsupported duration format with hours, minutes and seconds: %s", iso)
			}
		} else if strings.Contains(iso, "M") {
			// Format: PT#H#M (no seconds)
			_, err := fmt.Sscanf(iso, "PT%dH%dM", &hours, &minutes)
			if err != nil {
				return 0, fmt.Errorf("unsupported duration format with hours and minutes: %s", iso)
			}
		} else if strings.Contains(iso, "S") {
			// Format: PT#H#S (no minutes)
			_, err := fmt.Sscanf(iso, "PT%dH%dS", &hours, &seconds)
			if err != nil {
				return 0, fmt.Errorf("unsupported duration format with hours and seconds: %s", iso)
			}
		} else {
			// Format: PT#H only
			_, err := fmt.Sscanf(iso, "PT%dH", &hours)
			if err != nil {
				return 0, fmt.Errorf("unsupported duration format with hours only: %s", iso)
			}
		}
	} else if strings.Contains(iso, "M") && strings.Contains(iso, "S") {
		// Format: PT#M#S
		_, err := fmt.Sscanf(iso, "PT%dM%dS", &minutes, &seconds)
		if err != nil {
			return 0, fmt.Errorf("unsupported duration format with minutes and seconds: %s", iso)
		}
	} else if strings.Contains(iso, "M") {
		// Format: PT#M
		_, err := fmt.Sscanf(iso, "PT%dM", &minutes)
		if err != nil {
			return 0, fmt.Errorf("unsupported duration format with minutes only: %s", iso)
		}
	} else if strings.Contains(iso, "S") {
		// Format: PT#S
		_, err := fmt.Sscanf(iso, "PT%dS", &seconds)
		if err != nil {
			return 0, fmt.Errorf("unsupported duration format with seconds only: %s", iso)
		}
	} else {
		return 0, fmt.Errorf("unsupported duration format: %s", iso)
	}

	fmt.Println("hours, minutes, seconds", hours, minutes, seconds)
	totalMinutes := hours*60 + minutes

	// Option 1: Discard seconds entirely. For a video "PT4M8S", totalMinutes = 4.
	// Option 2: If you want to round up any partial minute, uncomment the following:
	// if seconds > 0 {
	//    totalMinutes++
	// }

	return totalMinutes, nil
}

// GetYouTubeVideoDuration fetches the video duration from the YouTube Data API.
func GetYouTubeVideoDuration(canonicalURL string) (int, error) {
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
	return dur, nil
}

// GetYouTubeVideoDurationMinutes fetches the video duration from the YouTube Data API
// and returns the duration in whole minutes (rounded up).
func CalculateCost(canonicalURL string) (int, error) {
	dur, err := GetYouTubeVideoDuration(canonicalURL)
	if err != nil {
		return 0, err
	}
	// round to nearest minute
	return dur * 25, nil
}

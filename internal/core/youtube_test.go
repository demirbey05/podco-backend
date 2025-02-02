package core_test

import (
	"strings"
	"testing"

	"github.com/demirbey05/auth-demo/internal/core"
	"github.com/joho/godotenv"
)

func TestYoutubeTitle(t *testing.T) {
	// Load environment variables same as main code
	if err := godotenv.Load("../../.env"); err != nil {
		t.Log("No .env file found")
	}
	title, err := core.GetYouTubeVideoTitle("https://www.youtube.com/watch?v=quIABSwc1Qc")
	if err != nil {
		t.Error(err)
	}
	if strings.Compare(title, "4 Hours of Ambient Study Music to Concentrate - Background Music For Concentration and Focus") != 0 {
		t.Error("Title mismatch")
	}
}

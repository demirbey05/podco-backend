package store

import (
	"context"
	"encoding/json"

	"github.com/demirbey05/auth-demo/db"
)

type Feedback []struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type FeedbackStore interface {
	InsertFeedback(ctx context.Context, userID string, feedback Feedback) (int, error)
}

type DBFeedbackStore struct {
	queries *db.Queries
}

func NewDBFeedbackStore(queries *db.Queries) *DBFeedbackStore {
	return &DBFeedbackStore{queries: queries}
}

func (s *DBFeedbackStore) InsertFeedback(ctx context.Context, userID string, feedback Feedback) error {

	// Convert feedback to JSON
	feedbackJSON, err := json.Marshal(feedback)
	if err != nil {
		return err
	}

	if err := s.queries.InsertFeedback(ctx, db.InsertFeedbackParams{
		CreatedBy: userID,
		Feedback:  feedbackJSON,
	}); err != nil {
		return err
	}
	return nil
}

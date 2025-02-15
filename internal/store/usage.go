package store

import (
	"context"

	"github.com/demirbey05/auth-demo/db"
)

type UsageStore interface {
	GetRemainingCredits(ctx context.Context, userID string) (int, error)
	DecrementCredit(ctx context.Context, userID string) (int, error)
}

type DBUsageStore struct {
	queries *db.Queries
}

func NewDBUsageStore(queries *db.Queries) *DBUsageStore {
	return &DBUsageStore{queries: queries}
}

func (s *DBUsageStore) GetRemainingCredits(ctx context.Context, userID string) (int, error) {
	remaining, err := s.queries.GetRemainingCredits(ctx, userID)
	if err.Error() == "no rows in result set" {
		return 15000, nil
	} else if err != nil {
		return 0, err
	}
	return int(remaining), nil
}

func (s *DBUsageStore) DecrementCredit(ctx context.Context, userID string) (int, error) {
	remainingExist, err := s.queries.IsCreditExist(ctx, userID)
	if err != nil {
		return 0, err
	}
	if !remainingExist {
		err := s.queries.InsertCredit(ctx, db.InsertCreditParams{UserID: userID, Credits: 14000})
		if err != nil {
			return 0, err
		}
		return 14000, nil
	}
	remaining, err := s.queries.DecrementCredit(ctx, db.DecrementCreditParams{UserID: userID, Credits: 1000})
	if err != nil {
		return 0, err
	}
	return int(remaining), nil
}

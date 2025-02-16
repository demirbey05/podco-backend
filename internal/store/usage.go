package store

import (
	"context"

	"github.com/demirbey05/auth-demo/db"
)

type UsageStore interface {
	GetRemainingCredits(ctx context.Context, userID string) (int, error)
	DecrementCredit(ctx context.Context, userID string, amount int) (int, error)
}

type DBUsageStore struct {
	queries *db.Queries
}

func NewDBUsageStore(queries *db.Queries) *DBUsageStore {
	return &DBUsageStore{queries: queries}
}

func (s *DBUsageStore) GetRemainingCredits(ctx context.Context, userID string) (int, error) {
	remaining, err := s.queries.GetRemainingCredits(ctx, userID)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 3000, nil
		}
		return 0, err
	}
	return int(remaining), nil
}

func (s *DBUsageStore) DecrementCredit(ctx context.Context, userID string, amount int) (int, error) {
	remainingExist, err := s.queries.IsCreditExist(ctx, userID)
	if err != nil {
		return 0, err
	}
	if !remainingExist {
		err := s.queries.InsertCredit(ctx, db.InsertCreditParams{UserID: userID, Credits: int32(3000 - amount)})
		if err != nil {
			return 0, err
		}
		return 14000, nil
	}
	remaining, err := s.queries.DecrementCredit(ctx, db.DecrementCreditParams{UserID: userID, Credits: int32(amount)})
	if err != nil {
		return 0, err
	}
	return int(remaining), nil
}

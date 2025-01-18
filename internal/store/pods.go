package store

import (
	"context"
	"time"

	"github.com/demirbey05/auth-demo/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type PodStore interface {
	GetPodsByLink(ctx context.Context, link string) ([]Pod, error)
	InsertPod(ctx context.Context, link string) (int, error)
	InsertArticle(ctx context.Context, podId int, userID, content string) error
	InsertQuiz(ctx context.Context, podId int, userID string) (int, error)
	InsertQuestion(ctx context.Context, quizId int, question string, options []string, correctIndex int) (int, error)
	InsertPodJob(ctx context.Context, podId int) (int, error)
	UpdatePodJob(ctx context.Context, jobId int, status int) error
}
type Pod struct {
	ID        int
	Link      string
	Title     string
	CreatedAt time.Time
}

type DBPodStore struct {
	queries *db.Queries
}

func NewDBPodStore(queries *db.Queries) *DBPodStore {
	return &DBPodStore{queries: queries}
}

func (s *DBPodStore) GetPodsByLink(ctx context.Context, link string) ([]Pod, error) {
	podDb, err := s.queries.GetPodByLink(ctx, link)
	if err != nil {
		return nil, err
	}
	pods := make([]Pod, len(podDb))
	for i, pod := range podDb {
		pods[i] = Pod{
			ID:        int(pod.ID),
			Link:      pod.Link,
			CreatedAt: pod.CreatedAt.Time,
		}
	}
	return pods, nil

}

// InsertPod inserts a new Pod and returns its ID.
func (s *DBPodStore) InsertPod(ctx context.Context, link string) (int, error) {
	pod, err := s.queries.InsertPod(ctx, link)
	if err != nil {
		return 0, err
	}
	return int(pod), nil
}

// InsertArticle inserts a new Article.
func (s *DBPodStore) InsertArticle(ctx context.Context, podId int, userID, content string) error {
	err := s.queries.InsertArticle(ctx, db.InsertArticleParams{
		PodID:       pgtype.Int4{Int32: int32(podId), Valid: true},
		CreatedBy:   userID,
		ArticleText: content,
	})
	return err
}

// InsertQuiz inserts a new Quiz and returns its ID.
func (s *DBPodStore) InsertQuiz(ctx context.Context, podId int, userID string) (int, error) {
	quiz, err := s.queries.InsertQuiz(ctx, db.InsertQuizParams{
		PodID:     pgtype.Int4{Int32: int32(podId), Valid: true},
		CreatedBy: userID,
	})
	if err != nil {
		return 0, err
	}
	return int(quiz), nil
}

// InsertQuestion inserts a new Question and returns its ID.
func (s *DBPodStore) InsertQuestion(ctx context.Context, quizId int, question string, options []string, correctIndex int) (int, error) {
	questionRecord, err := s.queries.InsertQuestion(ctx, db.InsertQuestionParams{
		QuizzesID:     pgtype.Int4{Int32: int32(quizId), Valid: true},
		QuestionText:  question,
		Options:       options,
		CorrectOption: int32(correctIndex),
	})
	if err != nil {
		return 0, err
	}
	return int(questionRecord), nil
}

func (s *DBPodStore) InsertPodJob(ctx context.Context, podId int) (int, error) {
	job, err := s.queries.InsertJob(ctx, int32(podId))
	if err != nil {
		return 0, err
	}
	return int(job), nil
}

func (s *DBPodStore) UpdatePodJob(ctx context.Context, jobId int, status int) error {
	return s.queries.UpdateJobStatusByID(ctx, db.UpdateJobStatusByIDParams{ID: int32(jobId), JobStatus: int32(status)})
}

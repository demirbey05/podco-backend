package store

import (
	"context"
	"fmt"
	"time"

	"github.com/demirbey05/auth-demo/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type PodStore interface {
	GetPodsByLink(ctx context.Context, link string) ([]Pod, error)
	InsertPod(ctx context.Context, link, title, userId string) (int, error)
	InsertArticle(ctx context.Context, podId int, content string) error
	InsertQuiz(ctx context.Context, podId int) (int, error)
	InsertQuestion(ctx context.Context, quizId int, question string, options []string, correctIndex int) (int, error)
	InsertPodJob(ctx context.Context, podId int) (int, error)
	UpdatePodJob(ctx context.Context, jobId int, status int) error
	GetArticleByPodID(ctx context.Context, podID int) (string, error)
	GetQuizByPodID(ctx context.Context, podID int) (QuizWithQuestions, error)
	GetJobStatus(ctx context.Context, jobID int) (int, error)
	GetPodsByUserID(ctx context.Context, userId string) ([]Pod, error)
	UpdatePodIsPublic(ctx context.Context, podID int, isPublic bool) error
	IsPodOwner(ctx context.Context, podID int, userID string) (bool, error)
}

type Pod struct {
	ID        int       `json:"id"`
	Link      string    `json:"link"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	IsPublic  bool      `json:"is_public"`
}

type QuizWithQuestions struct {
	ID        int        `json:"id"`
	PodID     int        `json:"pod_id"`
	Questions []Question `json:"questions"`
}

type Question struct {
	ID        int      `json:"id"`
	Text      string   `json:"question"`
	Options   []string `json:"options"`
	AnswerIdx int      `json:"correct_answer_index"`
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
			Title:     pod.Title,
			CreatedAt: pod.CreatedAt.Time,
			IsPublic:  pod.IsPublic.Bool,
		}
	}
	return pods, nil
}

func (s *DBPodStore) GetPodsByUserID(ctx context.Context, userId string) ([]Pod, error) {
	podDb, err := s.queries.GetPodsByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}
	pods := make([]Pod, len(podDb))
	for i, pod := range podDb {
		pods[i] = Pod{
			ID:        int(pod.ID),
			Link:      pod.Link,
			Title:     pod.Title,
			CreatedAt: pod.CreatedAt.Time,
			IsPublic:  pod.IsPublic.Bool,
		}
	}
	return pods, nil
}

// InsertPod inserts a new Pod and returns its ID.
func (s *DBPodStore) InsertPod(ctx context.Context, link, title, userId string) (int, error) {
	pod, err := s.queries.InsertPod(ctx, db.InsertPodParams{Link: link, Title: title, CreatedBy: userId})
	if err != nil {
		return 0, err
	}
	return int(pod), nil
}

// InsertArticle inserts a new Article.
func (s *DBPodStore) InsertArticle(ctx context.Context, podId int, content string) error {
	err := s.queries.InsertArticle(ctx, db.InsertArticleParams{
		PodID:       pgtype.Int4{Int32: int32(podId), Valid: true},
		ArticleText: content,
	})
	return err
}

// InsertQuiz inserts a new Quiz and returns its ID.
func (s *DBPodStore) InsertQuiz(ctx context.Context, podId int) (int, error) {
	quiz, err := s.queries.InsertQuiz(ctx, pgtype.Int4{Int32: int32(podId), Valid: true})
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

func (s *DBPodStore) GetArticleByPodID(ctx context.Context, podID int) (string, error) {
	article, err := s.queries.GetArticleByPodId(ctx, pgtype.Int4{Int32: int32(podID), Valid: true})
	if err != nil {
		return "", fmt.Errorf("error getting article: %w", err)
	}
	return article, nil
}

func (s *DBPodStore) GetQuizByPodID(ctx context.Context, podID int) (QuizWithQuestions, error) {
	quiz, err := s.queries.GetQuizByPodId(ctx, pgtype.Int4{Int32: int32(podID), Valid: true})
	if err != nil {
		return QuizWithQuestions{}, fmt.Errorf("error getting quiz: %w", err)
	}

	questions, err := s.queries.GetQuestionByQuizId(ctx, pgtype.Int4{Int32: int32(quiz.ID), Valid: true})
	if err != nil {
		return QuizWithQuestions{}, fmt.Errorf("error getting questions: %w", err)
	}

	result := QuizWithQuestions{
		ID:    int(quiz.ID),
		PodID: int(quiz.PodID.Int32),
	}

	for _, q := range questions {
		result.Questions = append(result.Questions, Question{
			ID:        int(q.ID),
			Text:      q.QuestionText,
			Options:   q.Options,
			AnswerIdx: int(q.CorrectOption),
		})
	}

	return result, nil
}

func (s *DBPodStore) GetJobStatus(ctx context.Context, jobID int) (int, error) {
	status, err := s.queries.GetJobStatusByID(ctx, int32(jobID))
	if err != nil {
		return 0, fmt.Errorf("error getting job status: %w", err)
	}
	return int(status), nil
}
func (s *DBPodStore) UpdatePodIsPublic(ctx context.Context, podID int, isPublic bool) error {
	return s.queries.UpdatePodIsPublic(ctx, db.UpdatePodIsPublicParams{ID: int32(podID), IsPublic: pgtype.Bool{Bool: isPublic, Valid: true}})
}
func (s *DBPodStore) IsPodOwner(ctx context.Context, podID int, userID string) (bool, error) {
	podInfo, err := s.queries.GetPodOwner(ctx, int32(podID))
	if err != nil {
		return false, fmt.Errorf("error getting pod owner: %w", err)
	}
	return podInfo.CreatedBy == userID || podInfo.IsPublic.Bool, nil
}

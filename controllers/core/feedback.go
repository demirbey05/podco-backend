package core

import (
	"github.com/demirbey05/auth-demo/db"
	"github.com/demirbey05/auth-demo/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insertFeedback(c *gin.Context, conn *pgxpool.Pool, queries *db.Queries) {
	userID := c.GetString("uuid")
	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var feedback store.Feedback
	if err := c.ShouldBindJSON(&feedback); err != nil {
		c.JSON(400, gin.H{"error": "invalid feedback"})
		return
	}

	if len(feedback) > 10 {
		c.JSON(400, gin.H{"error": "too many questions"})
		return
	}
	if len(feedback) == 0 {
		c.JSON(400, gin.H{"error": "no questions"})
		return
	}

	// Check character limit

	for _, question := range feedback {
		if len(question.Question) > 1000 {
			c.JSON(400, gin.H{"error": "question exceeds character limit"})
			return
		}
		if len(question.Answer) > 1000 {
			c.JSON(400, gin.H{"error": "answer exceeds character limit"})
			return
		}
	}

	feedbackStore := store.NewDBFeedbackStore(queries)

	if err := feedbackStore.InsertFeedback(c.Request.Context(), userID, feedback); err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"message": "feedback submitted successfully"})
}

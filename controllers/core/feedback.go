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

	feedbackStore := store.NewDBFeedbackStore(queries)

	if err := feedbackStore.InsertFeedback(c.Request.Context(), userID, feedback); err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"message": "feedback submitted successfully"})
}

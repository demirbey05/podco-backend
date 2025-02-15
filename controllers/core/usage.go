package core

import (
	"github.com/demirbey05/auth-demo/db"
	"github.com/demirbey05/auth-demo/internal/store"
	"github.com/gin-gonic/gin"
)

func getRemainingCredits(c *gin.Context, queries *db.Queries) {
	type resp struct {
		RemainingCredit int `json:"remaining_credit"`
	}
	userID := c.GetString("uuid")
	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	usageStore := store.NewDBUsageStore(queries)

	remainingCredit, err := usageStore.GetRemainingCredits(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, resp{RemainingCredit: remainingCredit})

}

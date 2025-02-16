package core

import (
	"fmt"

	"github.com/demirbey05/auth-demo/db"
	"github.com/demirbey05/auth-demo/internal/core"
	"github.com/demirbey05/auth-demo/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func createNewPod(c *gin.Context, conn *pgxpool.Pool, queries *db.Queries) {
	/* This endpoint will get youtube link and title from the request body and create a new pod in the database. */
	/* Then it schedules two job  : fetch the video transcription and send it to LLM to generate article  */
	/* After article is generated, it will be saved in the database */
	/* Then article will be sent to llm to generate quiz */

	var req struct {
		Link     string `json:"link" binding:"required"`
		Language string `json:"language" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": "bind error"})
		return
	}
	type resp struct {
		PodID           int `json:"pod_id"`
		JobId           int `json:"job_id"`
		RemainingCredit int `json:"remaining_credit"`
	}
	userID := c.GetString("uuid")
	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	tx, err := conn.Begin(c)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	defer tx.Rollback(c)

	qtx := queries.WithTx(tx)
	podStore := store.NewDBPodStore(qtx)
	usageStore := store.NewDBUsageStore(qtx)
	podID, jobID, remainingCredit, err := core.CreateNewPod(req.Link, userID, req.Language, podStore, usageStore)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	if err := tx.Commit(c); err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, resp{PodID: podID, JobId: jobID, RemainingCredit: remainingCredit})

}

func getPodsByUserID(c *gin.Context, conn *pgxpool.Pool, queries *db.Queries) {

	userID := c.GetString("uuid")
	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	podStore := store.NewDBPodStore(queries)

	pods, err := podStore.GetPodsByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"pods": pods})
}

func getArticle(c *gin.Context, conn *pgxpool.Pool, queries *db.Queries) {
	podID := c.Param("pod_id")
	var podIDInt int
	if _, err := fmt.Sscan(podID, &podIDInt); err != nil {
		fmt.Println(err)
		c.JSON(400, gin.H{"error": "invalid pod_id"})
		return
	}
	userID := c.GetString("uuid")
	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	podStore := store.NewDBPodStore(queries)

	isOwner, err := podStore.IsPodOwner(c.Request.Context(), podIDInt, userID)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	if !isOwner {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	article, err := podStore.GetArticleByPodID(c.Request.Context(), podIDInt)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"article": article})
}

func getQuiz(c *gin.Context, conn *pgxpool.Pool, queries *db.Queries) {
	podID := c.Param("pod_id")
	var podIDInt int
	if _, err := fmt.Sscan(podID, &podIDInt); err != nil {
		c.JSON(400, gin.H{"error": "invalid pod_id"})
		return
	}

	userID := c.GetString("uuid")
	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	podStore := store.NewDBPodStore(queries)
	isOwner, err := podStore.IsPodOwner(c.Request.Context(), podIDInt, userID)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	if !isOwner {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	quiz, err := podStore.GetQuizByPodID(c.Request.Context(), podIDInt)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"quiz": quiz})
}

func sharePod(c *gin.Context, conn *pgxpool.Pool, queries *db.Queries) {
	podIDStr := c.Param("pod_id")
	var podID int
	if _, err := fmt.Sscan(podIDStr, &podID); err != nil {
		c.JSON(400, gin.H{"error": "invalid pod_id"})
		return
	}

	userID := c.GetString("uuid")
	if userID == "" {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	podStore := store.NewDBPodStore(queries)

	// Only the pod owner can share (make public)
	isOwner, err := podStore.IsPodOwner(c.Request.Context(), podID, userID)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}
	if !isOwner {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	// Mark the pod as public
	if err := podStore.UpdatePodIsPublic(c.Request.Context(), podID, true); err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"message": "Pod is now public"})
}

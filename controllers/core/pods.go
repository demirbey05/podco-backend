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
		Link string `json:"link" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "bind error"})
		return
	}
	type resp struct {
		PodID int `json:"pod_id"`
		JobId int `json:"job_id"`
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

	podID, jobID, err := core.CreateNewPod(req.Link, podStore)
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

	c.JSON(200, resp{PodID: podID, JobId: jobID})

}

func getArticle(c *gin.Context, conn *pgxpool.Pool, queries *db.Queries) {
	podID := c.Param("pod_id")
	var podIDInt int
	if _, err := fmt.Sscan(podID, &podIDInt); err != nil {
		c.JSON(400, gin.H{"error": "invalid pod_id"})
		return
	}

	podStore := store.NewDBPodStore(queries)
	article, err := podStore.GetArticleByPodID(c.Request.Context(), podIDInt)
	if err != nil {
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

	podStore := store.NewDBPodStore(queries)
	quiz, err := podStore.GetQuizByPodID(c.Request.Context(), podIDInt)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"quiz": quiz})
}

func getJobStatus(c *gin.Context, conn *pgxpool.Pool, queries *db.Queries) {
	jobID := c.Param("job_id")
	var jobIDInt int
	if _, err := fmt.Sscan(jobID, &jobIDInt); err != nil {
		c.JSON(400, gin.H{"error": "invalid job_id"})
		return
	}

	podStore := store.NewDBPodStore(queries)
	status, err := podStore.GetJobStatus(c.Request.Context(), jobIDInt)
	if err != nil {
		c.JSON(500, gin.H{"error": "internal error"})
		return
	}

	c.JSON(200, gin.H{"status": status})
}

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

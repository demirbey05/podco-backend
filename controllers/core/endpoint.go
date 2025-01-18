package core

import (
	"github.com/demirbey05/auth-demo/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitCore(g *gin.Engine, conn *pgxpool.Pool, queries *db.Queries) {
	g.POST("/create-pod", func(ctx *gin.Context) {
		createNewPod(ctx, conn, queries)
	})
}

package core

import (
	"os"

	firebase "firebase.google.com/go"
	"github.com/demirbey05/auth-demo/controllers/middleware"
	"github.com/demirbey05/auth-demo/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitCore(g *gin.Engine, conn *pgxpool.Pool, queries *db.Queries, app *firebase.App) {
	// Configure CORS with FRONTEND_URL
	frontendURL := os.Getenv("FRONTEND_URL")
	config := cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	g.Use(cors.New(config))

	v1 := g.Group("/v1")
	protected := v1.Group("/protected")

	protected.Use(middleware.FirebaseAuthMiddleware(app))

	protected.POST("/try", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "You are authorized"})
	})

	v1.POST("/create-pod", func(ctx *gin.Context) {
		createNewPod(ctx, conn, queries)
	})
	v1.GET("/pods/:pod_id/article", func(ctx *gin.Context) {
		getArticle(ctx, conn, queries)
	})
	v1.GET("/pods/:pod_id/quiz", func(ctx *gin.Context) {
		getQuiz(ctx, conn, queries)
	})
	v1.GET("/jobs/:job_id/status", func(ctx *gin.Context) {
		getJobStatus(ctx, conn, queries)
	})
}

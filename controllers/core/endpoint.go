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

	protected.POST("/create-pod", func(ctx *gin.Context) {
		createNewPod(ctx, conn, queries)
	})
	protected.POST("/pods/share/:pod_id", func(ctx *gin.Context) {
		sharePod(ctx, conn, queries)
	})
	protected.GET("/pods/:pod_id/article", func(ctx *gin.Context) {
		getArticle(ctx, conn, queries)
	})
	protected.GET("/my-pods", func(ctx *gin.Context) {
		getPodsByUserID(ctx, conn, queries)
	})

	protected.GET("/pods/:pod_id/quiz", func(ctx *gin.Context) {
		getQuiz(ctx, conn, queries)
	})
}

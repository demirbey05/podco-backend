package main

import (
	"context"
	"os"

	"github.com/demirbey05/auth-demo/controllers/auth"
	"github.com/demirbey05/auth-demo/controllers/core"
	"github.com/demirbey05/auth-demo/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	url     string
	routers *gin.Engine
}

func NewServer() *Server {
	r := gin.Default()
	conn, queries, err := initStores(os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	addRoutes(r, conn, queries)

	url := os.Getenv("SERVICE_URL")

	return &Server{url: url, routers: r}
}
func (s *Server) Run() {
	s.routers.Run(s.url)
}
func initStores(postgresUrl string) (*pgxpool.Pool, *db.Queries, error) {

	conn, err := pgxpool.New(context.Background(), postgresUrl)
	if err != nil {
		return nil, nil, err
	}
	queries := db.New(conn)
	return conn, queries, nil

}

func addRoutes(r *gin.Engine, conn *pgxpool.Pool, queries *db.Queries) {
	auth.InitAuth(r)
	core.InitCore(r, conn, queries)
}

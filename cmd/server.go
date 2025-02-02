package main

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	"github.com/demirbey05/auth-demo/controllers/auth"
	"github.com/demirbey05/auth-demo/controllers/core"
	"github.com/demirbey05/auth-demo/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/option"
)

type Server struct {
	url     string
	routers *gin.Engine
	app     *firebase.App
}

func NewServer() *Server {
	r := gin.Default()
	conn, queries, err := initStores(os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	url := os.Getenv("SERVICE_URL")
	opt := option.WithCredentialsFile("./firebaseConfig.json")
	fireApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}
	addRoutes(r, conn, queries, fireApp)
	return &Server{url: url, routers: r, app: fireApp}
}
func (s *Server) Run() {
	// Run goose migrations

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

func addRoutes(r *gin.Engine, conn *pgxpool.Pool, queries *db.Queries, app *firebase.App) {
	auth.InitAuth(r)
	core.InitCore(r, conn, queries, app)
}

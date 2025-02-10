package main

import (
	"context"
	"os"

	firebase "firebase.google.com/go"
	"github.com/demirbey05/auth-demo/controllers/core"
	"github.com/demirbey05/auth-demo/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/option"
)

type Server struct {
	url     string
	conn    *pgxpool.Pool
	queries *db.Queries
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
	return &Server{url: url, routers: r, app: fireApp, conn: conn, queries: queries}
}
func (s *Server) Run() {
	// Add routers
	s.addRoutes()
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

func (s *Server) addRoutes() {
	core.InitCore(s.routers, s.conn, s.queries, s.app)
}

package main

import (
	"os"

	"github.com/demirbey05/auth-demo/internal/auth"
	"github.com/gin-gonic/gin"
)

type Server struct {
	url     string
	routers *gin.Engine
}

func NewServer() *Server {
	r := gin.Default()
	addRoutes(r)

	url := os.Getenv("SERVICE_URL")

	return &Server{url: url, routers: r}
}
func (s *Server) Run() {
	s.routers.Run(s.url)
}

func addRoutes(r *gin.Engine) {
	auth.InitAuth(r)
}

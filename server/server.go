package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	addr string
	db   *sql.DB
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() {
	router := gin.Default()

	apiRouter := router.Group("/api/v1")
	apiRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	log.Println("Server listening on", s.addr)
	if err := http.ListenAndServe(s.addr, router); err != nil {
		log.Fatalf("Error occured while listening on %s: %v", s.addr, err)
	}
}

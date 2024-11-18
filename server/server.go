package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/handler"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"
)

type Server struct {
	addr  string
	db    store.DataStore
	minio store.FileStore
}

func NewServer(addr string, db store.DataStore, minio store.FileStore) *Server {
	return &Server{
		addr:  addr,
		db:    db,
		minio: minio,
	}
}

// TODO: Implement the auth middleware later
func (s *Server) Run() {
	router := gin.Default()

	apiRouter := router.Group("/api/v1")
	apiRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	keyService := service.NewKeyService()
	service := service.NewService(s.db, s.minio, keyService)

	handler := handler.NewHandler(service)
	handler.RegisterRoutes(apiRouter)

	log.Println("Server listening on", s.addr)
	if err := http.ListenAndServe(s.addr, router); err != nil {
		log.Fatalf("Error occured while listening on %s: %v", s.addr, err)
	}
}

package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/handler"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/middleware"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"
)

type Server struct {
	addr  string
	db    store.DataStore
	minio store.FileStore
	redis store.CacheStore
	queue service.TaskQueue
}

func NewServer(addr string, db store.DataStore, minio store.FileStore, redis store.CacheStore, queue service.TaskQueue) *Server {
	return &Server{
		addr:  addr,
		db:    db,
		minio: minio,
		redis: redis,
		queue: queue,
	}
}

func (s *Server) Run() {
	router := gin.Default()

	unprotectedRouter := router.Group("/api/v1")
	unprotectedRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	authMiddleware := middleware.NewAuthMiddleware(s.db)
	protectedRouter := router.Group("/api/v1")
	protectedRouter.Use(authMiddleware.Middleware())

	keyService := service.NewKeyService()
	dummyFaceMatch := &service.DummyFaceMatchService{}
	dummyOcr := &service.DummyOcrService{}
	uuid := &service.UuidService{}

	service := service.NewService(s.db, s.minio, keyService, dummyFaceMatch, dummyOcr, s.queue, uuid)
	handler := handler.NewHandler(service)
	handler.RegisterRoutes(unprotectedRouter)
	handler.RegisterProtectedRoutes(protectedRouter)

	log.Println("Server listening on", s.addr)
	if err := http.ListenAndServe(s.addr, router); err != nil {
		log.Fatalf("Error occured while listening on %s: %v", s.addr, err)
	}
}

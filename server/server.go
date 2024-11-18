package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/config"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/handler"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"
	"github.com/minio/minio-go/v7"
)

type Server struct {
	addr  string
	db    store.DataStore
	minio *minio.Client
}

func NewServer(addr string, db store.DataStore, minio *minio.Client) *Server {
	return &Server{
		addr:  addr,
		db:    db,
		minio: minio,
	}
}

// TODO: Implement the auth middleware later
func (s *Server) Run() {
	router := gin.Default()

	// load configs
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error while config init: %v", err)
	}

	apiRouter := router.Group("/api/v1")
	apiRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	minioStore := service.NewMinioStore(s.minio, cfg.MinioBucket)
	keyService := service.NewKeyService()
	service := service.NewService(s.db, minioStore, keyService)

	handler := handler.NewHandler(service)
	handler.RegisterRoutes(apiRouter)

	log.Println("Server listening on", s.addr)
	if err := http.ListenAndServe(s.addr, router); err != nil {
		log.Fatalf("Error occured while listening on %s: %v", s.addr, err)
	}
}

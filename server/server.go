package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/handler"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/middleware"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/service"
	"github.com/justsushant/one2n-go-bootcamp/go-ekyc/store"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	addr  string
	db    store.DataStore
	minio store.FileStore
	redis store.CacheStore
	queue service.TaskQueue
}

type ServerConfig struct {
	Addr       string
	DataStore  store.DataStore
	FileStore  store.FileStore
	CacheStore store.CacheStore
	Queue      service.TaskQueue
}

func New(serverConfig *ServerConfig) *Server {
	return &Server{
		addr:  serverConfig.Addr,
		db:    serverConfig.DataStore,
		minio: serverConfig.FileStore,
		redis: serverConfig.CacheStore,
		queue: serverConfig.Queue,
	}
}

func (s *Server) Run() {
	router := gin.Default()

	// endpoint for swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	unprotectedRouter := router.Group("/api/v1")

	// endpoint for health check
	unprotectedRouter.GET("/health", HealthCheckHandler)

	authMiddleware := middleware.NewAuthMiddleware(s.db)
	protectedRouter := router.Group("/api/v1")
	protectedRouter.Use(authMiddleware.Middleware())

	keyService := service.NewKeyService()
	dummyFaceMatch := &service.DummyFaceMatchService{}
	dummyOcr := &service.DummyOcrService{}
	uuid := &service.UuidService{}

	serviceConfig := &service.ServiceConfig{
		DataStore:  s.db,
		FileStore:  s.minio,
		CacheStore: s.redis,
		KeyService: keyService,
		FaceMatch:  dummyFaceMatch,
		OCR:        dummyOcr,
		Queue:      s.queue,
		UUID:       uuid,
	}
	service := service.NewService(serviceConfig)
	handler := handler.NewHandler(service)
	handler.RegisterRoutes(unprotectedRouter)
	handler.RegisterProtectedRoutes(protectedRouter)

	log.Println("Server listening on", s.addr)
	if err := http.ListenAndServe(s.addr, router); err != nil {
		log.Fatalf("Error occured while listening on %s: %v", s.addr, err)
	}
}

// @Summary Health Check
// @Description Checks if the service is online
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string "Success Message"
// @Router /api/v1/health [get]
func HealthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "OK",
	})
}

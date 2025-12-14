package router

import (
	"context"
	"fmt"
	"workmgmt-api/internal/config"
	"workmgmt-api/internal/handlers"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	*gin.Engine
}

func NewRouter() *Router {
	// Set Gin mode based on environment
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()

	// Add middleware
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(corsMiddleware())

	return &Router{Engine: engine}
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RegisterHealthRoutes registers health check related routes
func RegisterHealthRoutes(r *Router, handler *handlers.HealthHandler) {
	v1 := r.Group("/api/v1")
	{
		health := v1.Group("/health")
		{
			health.GET("", handler.GetHealth)
			health.GET("/ready", handler.GetReadiness)
			health.GET("/live", handler.GetLiveness)
		}
	}
}

// RegisterAPIRoutes registers API specific routes
func RegisterAPIRoutes(r *Router, handler *handlers.APIHandler) {
	v1 := r.Group("/api/v1")
	{
		// API info
		v1.GET("/", handler.GetInfo)
		
		// /tickets routes
		tickets := v1.Group("/tickets")
		{
			tickets.POST("/", handler.CreateTicket)
			tickets.PUT("/:id", handler.UpdateTicket)
			tickets.GET("/:id", handler.GetTicket)
			tickets.POST("/:id/sync", handler.SyncTicket)
		}
		// /workflows routes
		workflows := v1.Group("/workflows")
		{
			workflows.GET("/", handler.ListWorkflows)
			workflows.POST("/", handler.CreateWorkflow)
		}
	}
}

// RegisterSwaggerRoutes registers Swagger documentation routes
func RegisterSwaggerRoutes(r *Router) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

type Server struct {
	httpServer *http.Server
	router     *Router
}

func NewServer(cfg *config.Config, router *Router) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%s", cfg.Port),
			Handler:      router.Engine,
			ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		},
		router: router,
	}
}

func (s *Server) Start() error {
	log.Printf("Server starting on port %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

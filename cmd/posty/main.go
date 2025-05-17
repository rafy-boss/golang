package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"syscall"
	"time"

	"github.com/azan-boss/posty/internal/config"
	"github.com/azan-boss/posty/internal/handler/middleware"
	"github.com/azan-boss/posty/internal/handler/websocket"
	"github.com/gin-contrib/cors"

	// "github.com/azan-boss/posty/internal/handler/websocket"
	"github.com/azan-boss/posty/internal/http/handler/chatroom"
	"github.com/azan-boss/posty/internal/http/handler/post"
	"github.com/azan-boss/posty/internal/http/handler/user"
	"github.com/azan-boss/posty/internal/storage/sqlite"
	"github.com/gin-gonic/gin"
)

func main() {
	slog.Info("application is starting..........")
	cfg := config.MustLoad()
	storage, err := sqlite.New(cfg)
	if err != nil {
		slog.Error("failed to create storage", "error", err)
		os.Exit(1)
	}

	// slog.Info("server address", "address", cfg.Env)
	router := gin.Default()
	
	// Add CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	
	protected := router.Group("/protected")
	protected.Use(middleware.AuthMiddleware())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World")
	})

	// Urls  will be list below here
	protected.POST("/post", post.New(storage))
	protected.POST("/chatroom", chatroom.CreatChatRoom(storage))
	router.POST("/register", user.New(storage))
	router.POST("/login", user.Login(storage))
	router.POST("/user", user.GetUser(storage))
	 router.GET("/ws/:roomId", GRWebsocket.HandleWebSocket(storage))
	 GRWebsocket.Init()
	 GRWebsocket.ConsumeQueue(storage)
	// router.GET("/ws", websocket.HandleWebSocket)

	// go websocket.HandleBroadcast()
	
	// Determine the address to listen on, prioritizing PORT env variable for Render
	var addr string
	
	// Check if PORT environment variable is set (for Render)
	port := os.Getenv("PORT")
	if port != "" {
		// For Render deployment, use the PORT environment variable
		addr = ":" + port
	} else {
		// For local development, use the address from config
		// If the address already has a colon prefix, use it as is
		if cfg.HttpServer.Address != "" {
			if cfg.HttpServer.Address[0] == ':' {
				addr = cfg.HttpServer.Address
			} else {
				// Add colon prefix if missing
				addr = ":" + cfg.HttpServer.Address
			}
		} else {
			// Default to port 8080 if no config is provided
			addr = ":8080"
		}
	}
	
	// Remove any double colons that might have been introduced
	if len(addr) >= 2 && addr[0] == ':' && addr[1] == ':' {
		addr = ":" + addr[2:]
	}
	
	slog.Info("server will start on", "address", addr)
	
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Error("server error", "error", err)
		}

	}()
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done
	slog.Info("shuting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shut down the server ")
	}
}

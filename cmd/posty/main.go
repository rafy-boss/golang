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
	server := &http.Server{
		Addr:    cfg.HttpServer.Address,
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

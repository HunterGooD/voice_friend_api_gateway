package main

import (
	"context"
	"errors"

	"github.com/HunterGooD/voice_friend_api_gateway/config"
	"github.com/HunterGooD/voice_friend_user_service/pkg/auth"
	"github.com/HunterGooD/voice_friend_user_service/pkg/logger"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.NewZapLogger()
	defer log.Sync()

	configPath := os.Getenv("CONFIG_PATH")
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Error("Error init config", err)
		panic(err)
	}

	tokenManager := auth.NewJWTGeneratorDefault("api_gateway")

	_, err = tokenManager.LoadPublicKeyFromFile(cfg.App.CertFilePath)
	if err != nil {
		log.Error("Error load public key from file", err)
		panic(err)
	}

	router := gin.Default()
	router.Use(gin.Recovery())

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server failed to start", err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTSTP)

	interrupt := <-quit

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	log.Info("Server is shutting down", map[string]any{
		"interrupt": interrupt,
	})

	if err := s.Shutdown(ctx); err != nil {
		log.Error("Server was unable to gracefully shutdown", err)
		panic(err)
	}

}

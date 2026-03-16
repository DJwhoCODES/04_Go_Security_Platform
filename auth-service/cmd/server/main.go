package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/djwhocodes/auth-service/internal/config"
	"github.com/djwhocodes/auth-service/internal/server"
	"github.com/djwhocodes/auth-service/pkg/logger"
)

func main() {

	cfg := config.LoadConfig()

	logger.Init(cfg.Log.Level)
	defer logger.Sync()

	srv := server.New(cfg)

	go srv.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv.Shutdown(ctx)
}

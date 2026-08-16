package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wyw14/cry005/internal/application"
	"github.com/wyw14/cry005/internal/config"
	"github.com/wyw14/cry005/internal/repository/memory"
	httptransport "github.com/wyw14/cry005/internal/transport/http"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	log, _ := zap.NewProduction()
	defer log.Sync()
	svc := application.New(memory.New())
	server := &http.Server{Addr: cfg.Address, Handler: httptransport.New(svc).Router(), ReadHeaderTimeout: cfg.RequestTimeout}
	go func() {
		log.Info("server started", zap.String("address", cfg.Address))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", zap.Error(err))
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}

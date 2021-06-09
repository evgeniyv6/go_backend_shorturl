package main

import (
	"context"
	"go.uber.org/zap"
	"go_backend_shorturl/configuration"
	"go_backend_shorturl/handler"
	"go_backend_shorturl/redisdb"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

var FSO configuration.OsFileSystem

func main() {
	logger := zap.NewExample()
	logger.Sync()

	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := configuration.ReadConfig("config.json", FSO)
	if err != nil {
		zap.S().Fatalf("Read configuration file error: %v", err)
	}

	srvAddress := config.Server.Address + ":" + config.Server.Port
	act, err := redisdb.NewPool(config.RedisDB.Address, config.RedisDB.Port)
	if err != nil {
		zap.S().Fatalf("Couldnot get redis connection error: %v", err)
	}
	defer func() {
		err := act.Close()
		if err != nil {
			zap.S().Errorf("Couldnot close redis connection.")
		}
	}()

	router := handler.NewGinRouter(config.Server.Protocol, srvAddress,act)

	srv := &http.Server{
		Addr:    srvAddress,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.S().Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	zap.S().Info("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout * time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.S().Fatalf("Server forced to shutdown: %s", err)
	}

	zap.S().Info("Server exiting")
}

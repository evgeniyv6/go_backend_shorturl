package main

import (
	"context"
	"net"

	//"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/configuration"
	"github.com/evgeniyv6/go_backend_shorturl/app/internal/handler"
	"github.com/evgeniyv6/go_backend_shorturl/app/internal/redisdb"

	"go.uber.org/zap"
)

var FSO configuration.OsFileSystem

func main() {
	logger := zap.NewExample()
	// Не может быть ошибки т.к. работаем с stdout
	//nolint: errcheck
	logger.Sync()

	// ReplaceGlobals replaces the global Logger and SugaredLogger, and returns a
	// function to restore the original values. It's safe for concurrent use.
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := configuration.ReadConfig("config.json", FSO)
	if err != nil {
		zap.S().Panicw("Read configuration file error.", "err", err)
	}

	srvAddress := net.JoinHostPort(config.Server.Address, config.Server.Port)
	zap.S().Infof("REDIS_ENDPOINT_URI = %s", os.Getenv("REDIS_ENDPOINT_URI"))
	zap.S().Infof("REDIS_ENDPOINT_PORT = %s", os.Getenv("REDIS_ENDPOINT_PORT"))
	act, err := redisdb.NewPool(os.Getenv("REDIS_ENDPOINT_URI"), os.Getenv("REDIS_ENDPOINT_PORT")) //(config.RedisDB.Address, config.RedisDB.Port)
	if err != nil {
		zap.S().Panicw("Couldnot get redis connection error.", "err", err)
	}
	defer func() {
		err := act.Close()
		if err != nil {
			zap.S().Error("Couldnot close redis connection.")
		}
	}()

	router := handler.NewGinRouter(config.Server.Protocol, srvAddress, act)

	err = router.Run()
	if err != nil {
		zap.S().Panicw("Start server failed.", "err", err)
	}
	zap.S().Info("srv started")
	//srv := &http.Server{
	//	Addr:    srvAddress,
	//	Handler: router,
	//}
	//
	//zap.S().Infof("Starting server. Listening address: %s on port: %s", config.Server.Address, config.Server.Port)
	//go func() {
	//	err := srv.ListenAndServe()
	//	zap.S().Infof("Server started. Listening address: %s on port: %s", config.Server.Address, config.Server.Port)
	//	if err != nil && err != http.ErrServerClosed {
	//		zap.S().Panicw("Start server error.", "err", err)
	//	}
	//}()

	// Listen for the interrupt signal
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify about shutdown
	stop()
	zap.S().Info("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to tell the server it has <timeout from config> seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout*time.Second)
	defer cancel()
	//if err := srv.Shutdown(ctx); err != nil {
	//	zap.S().Errorw("Server forced to shutdown.", "err", err)
	//}

	zap.S().Info("Server exiting.")
}

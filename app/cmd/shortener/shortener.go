package main

import (
	"context"
	"net"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/tracing"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/logger"

	"os/signal"
	"syscall"
	"time"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/configuration"
	"github.com/evgeniyv6/go_backend_shorturl/app/internal/handler"
	"github.com/evgeniyv6/go_backend_shorturl/app/internal/redisdb"

	_ "go.elastic.co/ecszap"
)

var FSO configuration.OsFileSystem

func main() {
	sugar := logger.NewZapWrapper()
	j := tracing.NewJaegerTracer("", sugar)
	tracer, closer := j.Init()
	defer closer.Close()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := configuration.ReadConfig("config.json", FSO)
	if err != nil {
		sugar.Panicw("Read configuration file error.", "err", err)
	}

	srvAddress := net.JoinHostPort(config.Server.Address, config.Server.Port)
	act, err := redisdb.NewPool(config.RedisDB.Address, config.RedisDB.Port, sugar, tracer)
	if err != nil {
		sugar.Panicw("Couldnot get redis connection error.", "err", err)
	}
	defer func() {
		err := act.Close()
		if err != nil {
			sugar.Error("Couldnot close redis connection.")
		}
	}()

	router := handler.NewGinRouter(config.Server.Protocol, srvAddress, act, sugar, tracer)

	err = router.Run()
	if err != nil {
		sugar.Panicw("Start server failed.", "err", err)
	}
	sugar.Info("srv started")

	// Listen for the interrupt signal
	<-ctx.Done()
	// Restore default behavior on the interrupt signal and notify about shutdown
	stop()
	sugar.Info("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to tell the server it has <timeout from config> seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout*time.Second)
	defer cancel()
	sugar.Info("Server exiting.")
}

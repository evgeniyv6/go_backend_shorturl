package main

import (
	"context"
	"net"
	"os"

	"go.elastic.co/ecszap"

	"os/signal"
	"syscall"
	"time"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/configuration"
	"github.com/evgeniyv6/go_backend_shorturl/app/internal/handler"
	"github.com/evgeniyv6/go_backend_shorturl/app/internal/redisdb"

	_ "go.elastic.co/ecszap"
	"go.uber.org/zap"
)

var FSO configuration.OsFileSystem

func main() {
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zap.New(core, zap.AddCaller())
	logger = logger.With(zap.String("app", "link cutter"))
	sugar := logger.Sugar()

	defer func() { _ = logger.Sync() }()

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := configuration.ReadConfig("config.json", FSO)
	if err != nil {
		sugar.Panicw("Read configuration file error.", "err", err)
	}

	srvAddress := net.JoinHostPort(config.Server.Address, config.Server.Port)
	act, err := redisdb.NewPool(config.RedisDB.Address, config.RedisDB.Port, sugar)
	if err != nil {
		sugar.Panicw("Couldnot get redis connection error.", "err", err)
	}
	defer func() {
		err := act.Close()
		if err != nil {
			sugar.Error("Couldnot close redis connection.")
		}
	}()

	router := handler.NewGinRouter(config.Server.Protocol, srvAddress, act, sugar)

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

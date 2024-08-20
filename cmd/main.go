package main

import (
	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal/adapter"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
	"github.com/hypay-id/backend-dashboard-hypay/internal/repository"
	"github.com/hypay-id/backend-dashboard-hypay/internal/server/http"
	"github.com/hypay-id/backend-dashboard-hypay/internal/server/http/controller"
	"github.com/hypay-id/backend-dashboard-hypay/internal/service"
	"go.uber.org/zap"
)

func main() {
	slog.NewLogger(slog.Info)
	// read from env
	envConfig, err := config.Reader()
	if err != nil {
		slog.Fatalw("failed to read config file", zap.Error(err))
	}

	// bind env to schema
	cfg := config.BindConfig(envConfig)

	// init repository for reads
	repoReads := repository.NewReadsRepo(cfg.Storage)

	// init repository for writes
	repoWrites := repository.NewWritesRepo(cfg.Storage)

	// adapter injector
	adptr := adapter.New(cfg.App)

	// init service/use-case/business logic
	svc := service.New(
		repoReads,
		repoWrites,
		cfg.App,
		adptr.MerchantCallback,
	)

	// http server will be used only for callback operation
	httpController := controller.NewController(cfg, svc.Transactions, svc.Merchants, svc.Users, svc.Providers)
	httpServer := http.NewHttpServer(cfg.HTTPServer, httpController)
	httpServer.ListenAndServe()
}

package http

import (
	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
	httpMiddleware "github.com/hypay-id/backend-dashboard-hypay/internal/server/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tylerb/graceful"
	"go.uber.org/zap"

	"github.com/hypay-id/backend-dashboard-hypay/internal/server/http/controller"
)

type ServerItf interface {
	ListenAndServe()
	Stop()
}

type Server struct {
	cfg            config.HTTPServer
	server         *graceful.Server
	echo           *echo.Echo
	httpController *controller.Controller
}

// New HTTP Server
func NewHttpServer(
	cfg config.HTTPServer,
	ctrl *controller.Controller,
) ServerItf {
	srv := &Server{
		echo:           echo.New(),
		cfg:            cfg,
		httpController: ctrl,
	}

	srv.connectCoreWithEcho()
	srv.initGracefulServer()

	return srv
}

func (s *Server) ListenAndServe() {
	// initialize framework

	// Setup HTTP Routes
	RegisterRoutes(s.echo, s.httpController)

	slog.Infow(
		"server_started",
		zap.String("address", s.echo.Server.Addr),
		zap.String("startup_type", "http"),
	)

	// serve http server gracefully
	_ = s.server.ListenAndServe()
}

func (s *Server) Stop() {
	gracefulTimeout := s.cfg.GracefulTimeout

	s.server.Stop(gracefulTimeout)
}

func (s *Server) connectCoreWithEcho() {

	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{"*"}, // Replace "*" with specific origins if needed
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}

	// Middleware
	s.echo.Use(httpMiddleware.RequestLogWithConfig())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.CORSWithConfig(corsConfig))

	RegisterRoutes(s.echo, s.httpController)
	setServerObj(s.echo, s.cfg)
}

func (h *Server) initGracefulServer() {
	e := h.echo
	h.server = &graceful.Server{
		Server:  e.Server,
		Timeout: h.cfg.GracefulTimeout,
		Logger:  graceful.DefaultLogger(),
	}
}

func setServerObj(e *echo.Echo, serverCfg config.HTTPServer) {
	e.Server.Addr = serverCfg.ListenAddress + ":" + serverCfg.Port
	if serverCfg.ReadTimeout > 0 {
		e.Server.ReadTimeout = serverCfg.ReadTimeout
	}
	if serverCfg.WriteTimeout > 0 {
		e.Server.WriteTimeout = serverCfg.WriteTimeout
	}
	if serverCfg.IdleTimeout > 0 {
		e.Server.IdleTimeout = serverCfg.IdleTimeout
	}
}

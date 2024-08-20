package middleware

import (
	"net"
	"os"
	"time"

	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	HeaderForwardedFor  = "x-forwarded-for"
	HeaderRealIP        = "x-real-ip"
	HealthCheckEndpoint = "/alive"
	ReadyCheckEndpoint  = "/ready"
)

func RequestLogWithConfig() echo.MiddlewareFunc {
	hostname, _ := os.Hostname()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err := next(c); err != nil {
				c.Error(err)
			}

			if req.URL.Path == HealthCheckEndpoint || req.URL.Path == ReadyCheckEndpoint {
				return nil
			}

			stop := time.Now()
			remoteAddr := req.RemoteAddr
			remoteIP, _, _ := net.SplitHostPort(remoteAddr)
			xffIP := req.Header.Get(HeaderForwardedFor)
			xRealIP := req.Header.Get(HeaderRealIP)

			slog.Infow("request_log",
				zap.String("hostname", hostname),
				zap.Int("status", res.Status),
				zap.String("remote_ip", remoteIP),
				zap.String("xff_ip", xffIP),
				zap.String("x_real_ip", xRealIP),
				zap.String("host", req.Host),
				zap.String("uri", req.RequestURI),
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.String("referer", req.Referer()),
				zap.Duration("latency", stop.Sub(start)),
				zap.String("latency_human", stop.Sub(start).String()),
			)

			return nil
		}
	}
}

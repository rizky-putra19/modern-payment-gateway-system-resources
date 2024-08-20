package controller

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/converter"
	"github.com/labstack/echo/v4"
)

type Controller struct {
	cfg                config.Schema
	transactionService internal.TransactionServiceItf
	merchantService    internal.MerchantServiceItf
	userService        internal.UserServiceItf
	providerService    internal.ProviderServiceItf
}

func NewController(
	cfg config.Schema,
	transaction internal.TransactionServiceItf,
	merchant internal.MerchantServiceItf,
	user internal.UserServiceItf,
	provider internal.ProviderServiceItf,
) *Controller {
	return &Controller{
		cfg:                cfg,
		transactionService: transaction,
		merchantService:    merchant,
		userService:        user,
		providerService:    provider,
	}
}

func (ctrl *Controller) ReturnOK(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (ctrl *Controller) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing Authorization header"})
		}

		token, err := request.ParseFromRequest(c.Request(), request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(ctrl.cfg.App.JWTSecret), nil
		})

		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}

		// Verify the token and check the expiration time
		if !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
		}

		claims := token.Claims.(jwt.MapClaims)
		currentTime := time.Now().Unix()
		if currentTime > converter.ToInt64(claims["exp"]) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "token has expired"})
		}

		c.Set("username", claims["username"])
		c.Set("userType", claims["userType"])
		c.Set("roleName", claims["roleName"])
		// slog.Infow("claimed token", "data", claims)

		return next(c)
	}
}

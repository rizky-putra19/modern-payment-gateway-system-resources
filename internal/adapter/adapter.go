package adapter

import (
	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal"
	merchantcallback "github.com/hypay-id/backend-dashboard-hypay/internal/adapter/merchant_callback"
)

type Adapter struct {
	MerchantCallback internal.MerchantCallbackItf
}

func New(cfg config.App) *Adapter {
	return &Adapter{
		MerchantCallback: merchantcallback.New(cfg),
	}
}

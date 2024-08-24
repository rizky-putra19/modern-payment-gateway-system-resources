package adapter

import (
	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal"
	"github.com/hypay-id/backend-dashboard-hypay/internal/adapter/jack"
	merchantcallback "github.com/hypay-id/backend-dashboard-hypay/internal/adapter/merchant_callback"
)

type Adapter struct {
	MerchantCallback internal.MerchantCallbackItf
	JackProvider     internal.JackProviderItf
}

func New(cfg config.App) *Adapter {
	return &Adapter{
		MerchantCallback: merchantcallback.New(cfg),
		JackProvider:     jack.New(cfg),
	}
}

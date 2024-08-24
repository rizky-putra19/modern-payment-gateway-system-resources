package service

import (
	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal"
	"github.com/hypay-id/backend-dashboard-hypay/internal/repository"
)

type Service struct {
	Transactions internal.TransactionServiceItf
	Merchants    internal.MerchantServiceItf
	Users        internal.UserServiceItf
	Providers    internal.ProviderServiceItf
}

func New(
	repoReads *repository.Repository,
	repoWrites *repository.Repository,
	cfg config.App,
	adptrMerchantCallback internal.MerchantCallbackItf,
	jackProvider internal.JackProviderItf,
) *Service {
	transactions := NewTransaction(
		repoReads.TransactionsReads,
		repoWrites.TransactionsWrites,
		repoReads.UserReads,
		repoReads.MerchantReads,
		repoWrites.MerchantWrites,
		cfg,
		jackProvider,
		repoReads.ProviderReads,
		repoWrites.ProviderWrites,
	)
	merchants := NewMerchant(repoReads.MerchantReads,
		repoWrites.MerchantWrites,
		repoReads.UserReads,
		adptrMerchantCallback,
		repoReads.TransactionsReads,
		repoReads.ProviderReads)
	providers := NewProvider(
		repoReads.TransactionsReads,
		repoReads.MerchantReads,
		repoReads.ProviderReads,
		repoReads.UserReads,
		cfg,
	)
	users := NewUser(repoReads.UserReads, repoWrites.UserWrites, cfg)

	return &Service{
		Transactions: transactions,
		Merchants:    merchants,
		Users:        users,
		Providers:    providers,
	}
}

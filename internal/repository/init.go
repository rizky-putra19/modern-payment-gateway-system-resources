package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
	"github.com/hypay-id/backend-dashboard-hypay/internal/repository/psql"
)

type Repository struct {
	db                 *sqlx.DB
	TransactionsReads  internal.TransactionsReadsRepositoryItf
	TransactionsWrites internal.TransactionsWritesRepositoryItf
	MerchantReads      internal.MerchantReadsRepositoryItf
	MerchantWrites     internal.MerchantWritesRepositoryItf
	ProviderReads      internal.ProviderReadsRepositoryItf
	ProviderWrites     internal.ProviderWritesRepositoryItf
	UserReads          internal.UserReadsRepositoryItf
	UserWrites         internal.UserWritesRepositoryItf
}

func NewReadsRepo(cfg config.Storage) *Repository {
	dbDriverReads, _ := initializeConnectionDbReads(cfg.PSQL["psqlReads"])
	slog.Infow("sql reads connection open", "dbname", cfg.PSQL["psqlReads"].DBName)

	transactionReads := psql.NewTransactionsReads(dbDriverReads)
	merchantReads := psql.NewMerchantReads(dbDriverReads)
	providerReads := psql.NewProviderReads(dbDriverReads)
	userReads := psql.NewUsersReads(dbDriverReads)

	return &Repository{
		db:                dbDriverReads,
		TransactionsReads: transactionReads,
		MerchantReads:     merchantReads,
		ProviderReads:     providerReads,
		UserReads:         userReads,
	}
}

func NewWritesRepo(cfg config.Storage) *Repository {
	dbDriverWrites, _ := initializeConnectionDbReads(cfg.PSQL["psqlWrites"])
	slog.Infow("sql writes connection open", "dbname", cfg.PSQL["psqlReads"].DBName)

	transactionWrites := psql.NewTransactionsWrites(dbDriverWrites)
	merchantWrites := psql.NewMerchantWrites(dbDriverWrites)
	userWrites := psql.NewUsersWrites(dbDriverWrites)
	providerWrites := psql.NewProviderWrites(dbDriverWrites)

	return &Repository{
		db:                 dbDriverWrites,
		TransactionsWrites: transactionWrites,
		MerchantWrites:     merchantWrites,
		UserWrites:         userWrites,
		ProviderWrites:     providerWrites,
	}
}

func initializeConnectionDbReads(psql config.PSQL) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		psql.Host,
		psql.Port,
		psql.User,
		psql.Password,
		psql.DBName,
	)
	db, err := sqlx.Open("postgres", connectionString)
	// defer db.Close()
	if err != nil {
		slog.Fatalw(
			"failed to initialize connection",
			"error",
			err.Error(),
		)
	}

	return db, err
}

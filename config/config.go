package config

import (
	"time"

	"github.com/spf13/viper"
)

// for testing purpose please change config path into absolute path
func Reader() (config EnvConfig, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// will replace with existing env
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

type Schema struct {
	HTTPServer HTTPServer
	Storage    Storage
	App        App
}

type App struct {
	JWTSecret           string
	OpsPassAdmin        string
	Type                string
	ProjectID           string
	PrivateKeyID        string
	PrivateKey          string
	ClientEmail         string
	ClientID            string
	AuthURI             string
	TokenURI            string
	AuthProviderCertURL string
	ClientCertURL       string
	Domain              string
	AppPassMail         string
}

type PSQL struct {
	DBName   string
	User     string
	Password string
	Host     string
	Port     int
}

type Storage struct {
	PSQL map[string]PSQL
}

type HTTPServer struct {
	ListenAddress   string
	Port            string
	GracefulTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

func BindConfig(env EnvConfig) Schema {
	return Schema{
		HTTPServer: HTTPServer{
			ListenAddress:   "0.0.0.0",
			Port:            env.ServerPort,
			GracefulTimeout: time.Second * 30,
			ReadTimeout:     time.Minute * 3,
			WriteTimeout:    time.Minute * 5,
			IdleTimeout:     time.Minute * 15,
		},
		Storage: Storage{
			PSQL: map[string]PSQL{
				"psqlWrites": {
					DBName:   env.StorageDatabaseNameWrites,
					User:     env.StorageDatabaseUsernameWrites,
					Password: env.StorageDatabasePasswordWrites,
					Host:     env.StorageDatabaseHostWrites,
					Port:     5432,
				},
				"psqlReads": {
					DBName:   env.StorageDatabaseNameReads,
					User:     env.StorageDatabaseUsernameReads,
					Password: env.StorageDatabasePasswordReads,
					Host:     env.StorageDatabaseHostReads,
					Port:     5432,
				},
			},
		},
		App: App{
			JWTSecret:           env.AppJWTSecret,
			OpsPassAdmin:        env.AppOpsPassAdmin,
			Type:                env.Type,
			ProjectID:           env.ProjectID,
			PrivateKeyID:        env.PrivateKeyID,
			PrivateKey:          env.PrivateKey,
			ClientEmail:         env.ClientEmail,
			ClientID:            env.ClientID,
			AuthURI:             env.AuthURI,
			TokenURI:            env.TokenURI,
			AuthProviderCertURL: env.AuthProviderCertURL,
			ClientCertURL:       env.ClientCertURL,
			Domain:              env.Domain,
			AppPassMail:         env.AppPassMail,
		},
	}
}

type EnvConfig struct {
	ServerPort                    string `mapstructure:"CONFIG_SERVER_PORT"`
	StorageDatabaseNameWrites     string `mapstructure:"CONFIG_STORAGE_DB_NAME_WRITES"`
	StorageDatabaseUsernameWrites string `mapstructure:"CONFIG_STORAGE_USERNAME_WRITES"`
	StorageDatabasePasswordWrites string `mapstructure:"CONFIG_STORAGE_PASSWORD_WRITES"`
	StorageDatabaseHostWrites     string `mapstructure:"CONFIG_STORAGE_HOST_WRITES"`
	StorageDatabaseNameReads      string `mapstructure:"CONFIG_STORAGE_DB_NAME_READS"`
	StorageDatabaseUsernameReads  string `mapstructure:"CONFIG_STORAGE_USERNAME_READS"`
	StorageDatabasePasswordReads  string `mapstructure:"CONFIG_STORAGE_PASSWORD_READS"`
	StorageDatabaseHostReads      string `mapstructure:"CONFIG_STORAGE_HOST_READS"`
	AppJWTSecret                  string `mapstructure:"CONFIG_APP_JWT_SECRET"`
	AppOpsPassAdmin               string `mapstructure:"CONFIG_OPERATIONS_PASSWORD"`
	Type                          string `mapstructure:"CONFIG_TYPE"`
	ProjectID                     string `mapstructure:"CONFIG_PROJECT_ID"`
	PrivateKeyID                  string `mapstructure:"CONFIG_PRIVATE_KEY_ID"`
	PrivateKey                    string `mapstructure:"CONFIG_PRIVATE_KEY"`
	ClientEmail                   string `mapstructure:"CONFIG_CLIENT_EMAIL"`
	ClientID                      string `mapstructure:"CONFIG_CLIENT_ID"`
	AuthURI                       string `mapstructure:"CONFIG_AUTH_URI"`
	TokenURI                      string `mapstructure:"CONFIG_TOKEN_URI"`
	AuthProviderCertURL           string `mapstructure:"CONFIG_AUTH_PROVIDER"`
	ClientCertURL                 string `mapstructure:"CONFIG_CLIENT"`
	Domain                        string `mapstructure:"CONFIG_DOMAIN"`
	AppPassMail                   string `mapstructure:"CONFIG_APP_MAIL"`
}

package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
	"google.golang.org/api/option"
)

func UploadFile(bucketName, fileName, objectName string, credentials []byte) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(credentials))
	if err != nil {
		return "", err
	}
	defer client.Close()

	f, err := os.Open(fileName)
	if err != nil {
		slog.Infof("os.Open: %v", err)
		return "", err
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// Membuat file dapat diakses secara publik
	if err := client.Bucket(bucketName).Object(objectName).ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	slog.Infof("File uploaded successfully: %v", fileName)
	return publicURL, nil
}

// Retrieve secret from Google Secret Manager
func GetSecret(config config.App) []byte {
	creds := dto.GoogleCredentials{
		Type:           config.Type,
		ProjectId:      config.ProjectID,
		PrivateKeyId:   config.PrivateKeyID,
		PrivateKey:     config.PrivateKey,
		ClientEmail:    config.ClientEmail,
		ClientId:       config.ClientID,
		AuthUri:        config.AuthURI,
		TokenUri:       config.TokenURI,
		AuthProvider:   config.AuthProviderCertURL,
		ClientCertUrl:  config.ClientCertURL,
		UniverseDomain: config.Domain,
	}

	credsJson, _ := json.Marshal(creds)

	return credsJson
}

package secret

import (
	"context"
	"os"
)

// EnvironmentSecretProvider Not safe for production environment. Use only when launching the application on local development environment.
type EnvironmentSecretProvider struct {
	databaseCredentials DatabaseCredentials
	hs256Secret         string
}

func NewEnvironmentSecretProvider() (*EnvironmentSecretProvider, error) {
	localSecretProvider := &EnvironmentSecretProvider{
		databaseCredentials: DatabaseCredentials{
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		},
		hs256Secret: os.Getenv("JWT_SIGNING_STRING"),
	}
	return localSecretProvider, nil
}

func (local *EnvironmentSecretProvider) GetDatabaseCredentials(ctx context.Context) (DatabaseCredentials, error) {
	return local.databaseCredentials, nil
}

func (local *EnvironmentSecretProvider) GetHS256Secret(ctx context.Context) (string, error) {
	return local.hs256Secret, nil
}

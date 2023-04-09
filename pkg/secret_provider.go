package secret

import (
	"context"
	"fmt"
)

type DatabaseCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Provider abstraction of secret management. Can be easily mock for test
type Provider interface {
	GetDatabaseCredentials(ctx context.Context) (DatabaseCredentials, error)
	GetHS256Secret(ctx context.Context) (string, error)
}

type Configuration struct {
	// true to use Vault to retreive secrets. Otherwise, retreive secret from environment (used for local development only)
	VaultSecretProvider bool

	// Correspond to the parameters https://developer.hashicorp.com/vault/api-docs/auth/kubernetes#role-1
	VaultRoleName string
}

func NewSecretProvider(ctx context.Context, conf Configuration) (Provider, error) {

	if conf.VaultSecretProvider {
		auth := KubernetesAuth{
			roleName: conf.VaultRoleName,
		}
		secretProvider, err := NewVaultSecretProviderWithKubernetesAuth(ctx, auth)
		if err != nil {
			return nil, fmt.Errorf("unable to initialize vault secret provider: %w", err)
		}
		return secretProvider, nil

	} else {
		secretProvider, err := NewEnvironmentSecretProvider()
		if err != nil {
			return nil, fmt.Errorf("unable to initialize local secret provider: %w", err)
		}
		return secretProvider, nil
	}

}

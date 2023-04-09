package secret

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	vault "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/api/auth/kubernetes"
)

type KubernetesAuth struct {
	roleName string
}

type VaultSecretProvider struct {
	client *vault.Client
}

func NewVaultSecretProviderWithKubernetesAuth(ctx context.Context, auth KubernetesAuth) (*VaultSecretProvider, error) {
	config := vault.DefaultConfig() // modify for more granular configuration
	log.Printf("connecting to vault @ %s", config.Address)

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize vault client: %w", err)
	}

	log.Printf("logging in to vault with kubernetes auth; role name: %s", auth.roleName)

	kubernetesAuth, err := kubernetes.NewKubernetesAuth(auth.roleName)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize kubernetes authentication method: %w", err)
	}

	authInfo, err := client.Auth().Login(ctx, kubernetesAuth)
	if err != nil {
		return nil, fmt.Errorf("unable to login using kubernetes auth method: %w", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no kubernetes info was returned after login")
	}

	log.Println("connecting to vault: success!")

	vaultSecretProvider := &VaultSecretProvider{
		client: client,
	}
	return vaultSecretProvider, nil

}

func (vault *VaultSecretProvider) GetDatabaseCredentials(ctx context.Context) (DatabaseCredentials, error) {
	lease, err := vault.client.Logical().ReadWithContext(ctx, "database/creds/innsecure")
	if err != nil {
		return DatabaseCredentials{}, fmt.Errorf("unable to read secret: %w", err)
	}

	b, err := json.Marshal(lease.Data)
	if err != nil {
		return DatabaseCredentials{}, fmt.Errorf("malformed credentials returned: %w", err)
	}

	var credentials DatabaseCredentials

	if err := json.Unmarshal(b, &credentials); err != nil {
		return DatabaseCredentials{}, fmt.Errorf("unable to unmarshal credentials: %w", err)
	}

	log.Println("getting temporary database credentials from vault: success!")

	// raw secret is included to renew database credentials
	return credentials, nil
}

func (vault *VaultSecretProvider) GetHS256Secret(ctx context.Context) (string, error) {
	secret, err := vault.client.KVv2("secret").Get(context.Background(), "innsecure")
	if err != nil {
		return "", fmt.Errorf("unable to read secret: %w", err)
	}

	// data map can contain more than one key-value pair,
	// in this case we're just grabbing one of them
	value, ok := secret.Data["hs256"].(string)
	if !ok {
		return "", fmt.Errorf("unable to read secret: %w", err)
	}

	return value, nil
}

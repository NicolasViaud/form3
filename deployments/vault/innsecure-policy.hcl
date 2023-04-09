path "secret/data/innsecure" {
  capabilities = ["read"]
}

# Allows read-only access to the secret path that will be used
# by Vault to handle generation of dynamic database credentials.
path "database/creds/innsecure" {
  capabilities = ["read"]
}
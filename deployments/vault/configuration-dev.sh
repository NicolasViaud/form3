#!/bin/sh

###############################################################################################
##               *** WARNING - INSECURE - DO NOT USE IN PRODUCTION ***                       ##
## This script is to simulate operations a Vault operator would perform and, as such,        ##
## is not a representation of best practices in production environments.                     ##
## https://learn.hashicorp.com/tutorials/vault/pattern-approle?in=vault/recommended-patterns ##
###############################################################################################


#####################################
########## ACCESS POLICIES ##########
#####################################

# Add policies for the various roles we'll be using
# ref: https://www.vaultproject.io/docs/concepts/policies
vault policy write innsecure-policy /vault/config/innsecure-policy.hcl


#####################################
###### KUBERNETES AUTH METHOD #######
#####################################

# Enable Kubernetes auth method utilized by our application
# ref: https://developer.hashicorp.com/vault/docs/auth/kubernetes
vault auth enable kubernetes

vault write auth/kubernetes/config \
    kubernetes_host="https://$KUBERNETES_PORT_443_TCP_ADDR:443"

vault write auth/kubernetes/role/innsecure \
        bound_service_account_names=innsecure \
        bound_service_account_namespaces=innsecure \
        policies=innsecure-policy \
        ttl=24h


#####################################
########## STATIC SECRETS ###########
#####################################

# Seed the secret store
vault kv put "secret/innsecure" "hs256=84GIuIXeTTuf3ztA2werUAr5uk30ICF0VJWvTP8xAzBw0hiY5VgnefQ5SsW7vr1I"

#####################################
########## DYNAMIC SECRETS ##########
#####################################

# Enable the database secrets engine
# ref: https://www.vaultproject.io/docs/secrets/databases
vault secrets enable database

# Configure Vault's connection to our db, in this case PostgreSQL
# ref: https://www.vaultproject.io/api/secret/databases/postgresql
vault write database/config/postgres \
    plugin_name=postgresql-database-plugin \
    allowed_roles="innsecure" \
    connection_url="postgresql://{{username}}:{{password}}@postgres.innsecure.svc.cluster.local/innsecure?sslmode=disable" \
    username="root" \
    password="YWRtaW4xMjM0"


# Allow Vault to create roles dynamically. For test purpose, the time to live is 1 day
vault write database/roles/innsecure \
    db_name=postgres \
    creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA Public TO \"{{name}}\" ;" \
    renew_statements="ALTER ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA Public TO \"{{name}}\" ;"

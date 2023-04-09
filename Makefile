.DEFAULT_GOAL := default

clean:
	@rm ./bin/*

default: install-deps build test

build: install-deps
	@./bin/hey go build -o bin/innsecure ./cmd/innsecure
	@./bin/hey go build -o bin/token ./cmd/token

bin_dir:
	@mkdir -p ./bin

install-deps: install-hey install-goimports

install-hey: bin_dir
	@curl -L --insecure https://github.com/rossmcf/hey/releases/download/v1.0.0/installer.sh | bash
	@mv hey bin

install-goimports:
	@if [ ! -f ./goimports ]; then \
		cd ~ && go install golang.org/x/tools/cmd/goimports@latest; \
	fi

install-kind:
	@go install sigs.k8s.io/kind@v0.13.0

test:
	@echo "executing tests..."
	go test github.com/form3tech/innsecure

# package for release to candidates (ignore for test exercise)
package-%:
	@echo $*
	@cd ..&& pwd && tar -czvf innsecure-$*.tar.gz --exclude={".git",".github","bin","releases"} innsecure
	@mkdir -p releases
	@mv ../innsecure-$*.tar.gz releases 

get-docker-images:
	@docker build . -t form3/innsecure
	@docker pull postgres:12

deploy-vault:
	@helm repo add hashicorp https://helm.releases.hashicorp.com
	@helm repo update
	@kubectl create namespace vault
	@kubectl create configmap config -n vault --from-file=deployments/vault/innsecure-policy.hcl --from-file=deployments/vault/configuration-dev.sh
	@helm install vault hashicorp/vault -n vault --values deployments/vault/helm-vault-dev.yml --wait
	@kubectl wait --for=condition=Ready pod -n vault vault-0
	@kubectl wait --for=condition=Ready pod -n innsecure -l app=postgres
	@kubectl exec -n vault vault-0 -- sh /vault/config/configuration-dev.sh

undeploy-vault:
	@helm uninstall vault -n vault
	@kubectl delete namespace vault

deploy-application:
	@kubectl apply -f ./k8s

undeploy-application:
	@kubectl delete -f ./k8s

start-kind: get-docker-images
	@kind create cluster
	@kind load docker-image form3/innsecure
	@kind load docker-image postgres:12
	
stop-kind:
	@kind delete cluster
	
start-local: install-kind start-kind deploy-application deploy-vault 

stop-local: stop-kind
	@docker image rm form3/innsecure

.PHONY: clean build test package-% start-local stop-local

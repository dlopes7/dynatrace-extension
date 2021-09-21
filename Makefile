VERSION ?= 0.0.7
IMG=quay.io/dlopes7/dt-extension

docker-build: ## Build docker image with the manager.
	docker build -t ${IMG}:${VERSION} -t ${IMG}:latest .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}

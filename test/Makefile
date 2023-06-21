repo=k3d-local-registry:12345

build:
	docker build -t $(repo)/hello:latest .
	-docker push $(repo)/hello:latest

deploy-blue:
	helm upgrade --install example-blue chart --set image.repository=$(repo)/hello --set image.tag=latest

deploy-green:
	helm upgrade --install example-green chart --set image.repository=$(repo)/hello --set image.tag=latest --set extraSelectorLabels.version=green

deploy-nats:
	helm upgrade --install nats nats --repo https://nats-io.github.io/k8s/helm/charts --version 0.17.5
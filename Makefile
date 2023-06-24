
repo=k3d-local-registry:12345
cluster?="test"
https-port?=8443
http-port?=8080

k3d:
	-k3d registry create local-registry --port 12345
	# https://github.com/k3d-io/k3d/issues/209
	@export K3D_FIX_DNS=1 && k3d cluster create $(cluster) --registry-use k3d-local-registry:12345 -p "$(https-port):443@loadbalancer" -p "$(http-port):80@loadbalancer" --k3s-arg "--disable=traefik@server:0"
	@kubectl config use-context k3d-$(cluster)
	@helm install ingress-nginx --repo https://kubernetes.github.io/ingress-nginx ingress-nginx
	@helm install nats nats --repo https://nats-io.github.io/k8s/helm/charts --version 0.17.5

build:
	docker build -t $(repo)/switchery:latest .
	-docker push $(repo)/switchery:latest

deploy:
	helm upgrade --install switchery chart -f example.yaml

stress:
	k6 run k6-test.js --insecure-skip-tls-verify

blue:
	kubectl exec deployment/nats-box -- nats pub default/example blue

green:
	kubectl exec deployment/nats-box -- nats pub default/example green

switch:
	kubectl exec deployment/nats-box -- nats pub default/example switch
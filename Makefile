
repo=k3d-local-registry:12345

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
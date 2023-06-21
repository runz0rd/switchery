
repo=k3d-local-registry:12345

build:
	docker build -t $(repo)/switchy:latest .
	-docker push $(repo)/switchy:latest

deploy:
	helm upgrade --install switchy test/chart -f values.yaml

stress:
	# kubectl port-forward svc/example 8080:80
	k6 run k6-test.js

blue:
	kubectl patch svc/example -p '{"spec":{"selector":{"version": "blue"}}}'

green:
	kubectl patch svc/example -p '{"spec":{"selector":{"version": "green"}}}'

switch:
	kubectl exec deployment/nats-box -- nats pub foo default/example
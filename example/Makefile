repo=k3d-local-registry:12345

build:
	docker build -t $(repo)/hello:latest .
	-docker push $(repo)/hello:latest

deploy-blue:
	helm upgrade --install example-blue chart --set image.repository=$(repo)/hello --set image.tag=latest -f blue.yaml

deploy-green:
	helm upgrade --install example-green chart --set image.repository=$(repo)/hello --set image.tag=latest -f green.yaml

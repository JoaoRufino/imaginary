OK_COLOR=\033[32;01m
NO_COLOR=\033[0m

build: swagger
	@echo "$(OK_COLOR)==> Compiling binary$(NO_COLOR)"
	go test && go build -o bin/imaginary

test:
	go test

install:
	go get -u .

benchmark: build
	bash benchmark.sh

docker-build:
	@echo "$(OK_COLOR)==> Building Docker image$(NO_COLOR)"
	docker build --no-cache=true --build-arg IMAGINARY_VERSION=$(VERSION) -t h2non/imaginary:$(VERSION) .

docker-push:
	@echo "$(OK_COLOR)==> Pushing Docker image v$(VERSION) $(NO_COLOR)"
	docker push h2non/imaginary:$(VERSION)

docker: docker-build docker-push

check-swagger:
	which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	GO111MODULE=on go mod vendor && GO111MODULE=off swagger generate spec -o public_html/api-docs/swagger.json --scan-models

serve-swagger: check-swagger
	swagger serve -F=swagger public_html/api-docs/swagger.json

.PHONY: test benchmark docker-build docker-push docker

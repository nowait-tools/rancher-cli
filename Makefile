.PHONY: build push test

PWD := `pwd`
IMAGE_NAME := rancher-cli
IMAGE_TAG := 0.3.0
BIN := rancher-cli
PACKAGE_PATH := github.com/nowait-tools/rancher-cli

build:
	docker run -it \
	    -v $(PWD):/go/src/github.com/nowait-tools/rancher-cli \
	    -e "GOOS=linux" \
	    -e "GOARCH=amd64" \
	    -w /go/src/$(PACKAGE_PATH)  golang:1.7-alpine \
        go build -o $(BIN)
	docker build -t nowait/$(IMAGE_NAME):$(IMAGE_TAG) .

push:
	docker push nowait/$(IMAGE_NAME):$(IMAGE_TAG)

test:
	docker run -it \
	    -v $(PWD):/go/src/github.com/nowait-tools/rancher-cli \
	    -e "GOOS=linux" \
	    -e "GOARCH=amd64" \
	    -w /go/src/$(PACKAGE_PATH)  golang:1.7-alpine \
        go test ./...

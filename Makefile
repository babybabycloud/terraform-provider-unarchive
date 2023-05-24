UNARCHIVE_PATH := internal/unarchive

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build: 
	go build -o ${HOME}/workspace/terraform-plugin/

.PHONY: generate-doc
generate-doc:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run

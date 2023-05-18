UNARCHIVE_PATH := internal/unarchive

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build: 
	go build -o ${HOME}/workspace/terraform-plugin/

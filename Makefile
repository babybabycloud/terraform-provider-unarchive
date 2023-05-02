.PHONY: test
test: 
	go test -v ./...

.PHONY: build
build: 
	go build -o /home/vagrant/workspace/terraform-plugin/

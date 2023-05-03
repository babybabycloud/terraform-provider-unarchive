.PHONY: test
test: clean-test pre-test
	go test -v ./...

.PHONY: build
build: 
	go build -o /home/vagrant/workspace/terraform-plugin/

.PHONY: clean-test
clean-test:
	rm -rf unarchive/vali-helper-master unarchive/master.zip

.PHONY: pre-test
pre-test:
	wget -O unarchive/master.zip https://gitee.com/babybabycloud/vali-helper/repository/archive/master.zip

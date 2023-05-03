UNARCHIVE_PATH := internal/unarchive

.PHONY: test
test: clean-test pre-test
	go test -v ./...

.PHONY: build
build: 
	go build -o /home/vagrant/workspace/terraform-plugin/

.PHONY: clean-test
clean-test:
	rm -rf $(UNARCHIVE_PATH)/vali-helper-master $(UNARCHIVE_PATH)/master.zip

.PHONY: pre-test
pre-test:
	wget -O $(UNARCHIVE_PATH)/master.zip https://gitee.com/babybabycloud/vali-helper/repository/archive/master.zip

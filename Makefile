default: build

clean:
	$(RM) ./bin/docker-machine-driver-ucloud
	$(RM) $(GOPATH)/bin/docker-machine-driver-ucloud

build: clean
	GOGC=off go build -i -o ./bin/docker-machine-driver-ucloud ./bin

install: build
	cp ./bin/docker-machine-driver-ucloud $(GOPATH)/bin/

.PHONY: build install

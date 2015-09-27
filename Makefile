default: build

clean:
	$(RM) ./bin/docker-machine-ucloud
	$(RM) $(GOPATH)/bin/docker-machine-ucloud

build: clean
	GOGC=off go build -i -o ./bin/docker-machine-ucloud ./bin

install: build
	cp ./bin/docker-machine-ucloud $(GOPATH)/bin/

.PHONY: build install
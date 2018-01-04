CONTAINER_NAME=package-server
SERVER_DIR_NAME=reposerver

.PHONY: build
.PHONY: fmt
.PHONY: vet
.PHONY: docker.build
.PHONY: test.unit
.PHONY: test
.PHONY: clean

build: fmt vet reposerver

fmt:
	go $@ ./...

vet:
	go $@ ./...

reposerver:
	go build -o $(SERVER_DIR_NAME) cmd/$(SERVER_DIR_NAME)/reposerver.go

client:
	cd cmd/$(CLIENT_DIR_NAME); go build -o client && cp client ../..

build.docker:
	docker build -t $(CONTAINER_NAME) .

test: test.unit 

test.unit: fmt vet
	go test -cover ./...

clean:
	-rm reposerver
	-rm client
	-rm packagetree.zip

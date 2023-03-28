# TODO: change project name
PROJECT = test-proxy
MAIN = cmd/main.go
GIT= github.com/nnickie23/test_proxy
# REGISTRY=registry.wtotem/webtotem

# build - target for building go application
build: $(MAIN)
	go build -v -ldflags '-w -s' -o bin/$(PROJECT) $(MAIN)

# docker - target for building docker image with go application
# Need WT_TAG environment variable
docker:
ifndef WT_TAG
	$(error WT_TAG is undefined)
endif
	docker build -t $(REGISTRY)/$(PROJECT):$(WT_TAG) -f build/Dockerfile .

lint:
	golangci-lint run -E revive -E gosec -E stylecheck -E bodyclose -E cyclop -E gofmt -D structcheck

test:
	go test -v ./...
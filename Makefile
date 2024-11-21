.PHONY: test build

MODULE := github.com/Scorpio69t/gcloc
APP_NAME := gcloc
VERSION_FULL := $(shell git describe --tags --always)
VERSION := $(shell echo $(VERSION_FULL) | awk -F'-' '{print $1}')
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date +"%Y-%m-%dT%H:%M:%SZ")

DOCKER_REPO := scorpio69t/$(APP_NAME)

LD_FLAGS := -ldflags "-X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.GitCommit=$(COMMIT) -X $(MODULE)/cmd.BuildDate=$(DATE)"

build: cleanup-package
	mkdir -p bin
	GO111MODULE=on go build $(LD_FLAGS) -o ./bin/gcloc app/gcloc/main.go

build-docker: cleanup-package
	docker build \
		--build-arg VERSION=$(VERSION) \
    	--build-arg GIT_COMMIT=$(COMMIT) \
    	--build-arg BUILD_DATE=$(DATE) \
    	-t $(DOCKER_REPO):$(VERSION) .

push-docker:
	docker push $(DOCKER_REPO):$(VERSION)

update-package:
	GO111MODULE=on go get -u github.com/Scorpio69t/gcloc

cleanup-package:
	GO111MODULE=on go mod tidy

gcloc-linux-amd64: cleanup-package
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build $(LD_FLAGS) -o ./bin/gcloc-linux-amd64 app/gcloc/main.go

gcloc-linux-arm64: cleanup-package
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GO111MODULE=on go build $(LD_FLAGS) -o ./bin/gcloc-linux-arm64 app/gcloc/main.go

gcloc-darwin-amd64: cleanup-package
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on go build $(LD_FLAGS) -o ./bin/gcloc-darwin-amd64 app/gcloc/main.go

gcloc-darwin-arm64: cleanup-package
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 GO111MODULE=on go build $(LD_FLAGS) -o ./bin/gcloc-darwin-arm64 app/gcloc/main.go

gcloc-windows-amd64: cleanup-package
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GO111MODULE=on go build $(LD_FLAGS) -o ./bin/gcloc-windows-amd64.exe app/gcloc/main.go

#run-example:
#	GO111MODULE=on go run examples/languages/main.go
#	GO111MODULE=on go run examples/files/main.go

test:
	GO111MODULE=on go test -v

test-cover:
	GO111MODULE=on go test -v -coverprofile=coverage.out

clean:
	rm -rf bin

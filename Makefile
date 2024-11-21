.PHONY: test build

MODULE := github.com/Scorpio69t/gcloc
VERSION := $(shell git describe --tags --always)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date +"%Y-%m-%dT%H:%M:%SZ")

build:
	mkdir -p bin
	GO111MODULE=on go build -ldflags "-X $(MODULE)/cmd.Version=$(VERSION) -X $(MODULE)/cmd.GitCommit=$(COMMIT) -X $(MODULE)/cmd.BuildDate=$(DATE)" -o ./bin/gcloc app/gcloc/main.go

update-package:
	GO111MODULE=on go get -u github.com/Scorpio69t/gcloc

cleanup-package:
	GO111MODULE=on go mod tidy

gcloc-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o ./bin/gcloc-linux-amd64 app/gcloc/main.go

gcloc-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GO111MODULE=on go build -o ./bin/gcloc-linux-arm64 app/gcloc/main.go

gcloc-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 GO111MODULE=on go build -o ./bin/gcloc-darwin-amd64 app/gcloc/main.go

gcloc-darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 GO111MODULE=on go build -o ./bin/gcloc-darwin-arm64 app/gcloc/main.go

gcloc-windows-amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GO111MODULE=on go build -o ./bin/gcloc-windows-amd64.exe app/gcloc/main.go

#run-example:
#	GO111MODULE=on go run examples/languages/main.go
#	GO111MODULE=on go run examples/files/main.go

test:
	GO111MODULE=on go test -v

test-cover:
	GO111MODULE=on go test -v -coverprofile=coverage.out

PROJECTNAME=$(shell echo $(shell pwd) | sed 's!.*/!!')

# Go related variables.
BASE=$(shell pwd)
GOBASE=$(BASE)/api
GOPATH=$(GOBASE)/vendor
GOBIN=$(BASE)/bin
GOBUILD=$(GOBIN)/$(PROJECTNAME)
GOFILES=$(wildcard *.go)

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID=/tmp/.$(PROJECTNAME).pid

## clean: Clean build files
clean:
	@echo "  >  Cleaning build cache"
	@rm -rf $(GOBIN)

## generate: run go generate
generate:
	@go generate ./...

## build: Build API Binary for Host OS
build: clean
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o $(GOBUILD) $(GOBASE)

build-linux: clean
	@echo "  >  Building binary..."
	@GOOS=linux CGO_ENABLED=0 go build -o $(GOBUILD) $(GOBASE)

## docker: Tag new docker image based on build
docker: build-linux
	@docker build -t quirk-api .

## start: Start the API in dev mode
start:
	@CONFIG=$(GOBASE)/config.toml $(GOBUILD)

## run: Build and start the API
run: build start

## End to End testing
test-e2e:
	go test ./api/tests -v

## Unit and integration tests
test:
	go test `go list ./api/... | grep -v tests` -v -cover

test-all:
	go test ./api/... -v

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
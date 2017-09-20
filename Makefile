PROJECT := gcd
BIN_DIR := $(GOPATH)/bin
RELEASE_DIR := ./release

VERSION ?= v0.0.1
PLATFORM ?= linux
ARCH ?= amd64
RELEASE_PATH := $(RELEASE_DIR)/$(PROJECT)

PATH_MAIN_PACKAGE := ./
PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test lint build clear

build: clear $(RELEASE_PATH)

build-docker:
	@echo "---> Building the project using Dockerfile"
	@docker build -t guiferpa/$(PROJECT):$(VERSION) .

clear:
	@echo "---> Cleaning up directory"
	@rm -rf $(RELEASE_DIR)

$(RELEASE_PATH):
	@echo "---> Building the project"
	@mkdir -p $(RELEASE_DIR)
	@CGO_ENABLED=0 GOOS=$(PLATFORM) GOARCH=$(ARCH) go build -o $(RELEASE_PATH) $(PATH_MAIN_PACKAGE)

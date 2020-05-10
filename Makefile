VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

CURRENT_VERSION := $(shell bin/current-version)

ifndef CURRENT_VERSION
	CURRENT_VERSION := 0.0.0
endif

NEXT_VERSION := $(shell semver -c -i $(RELEASE_TYPE) $(CURRENT_VERSION))
DOCKER_NEXT_VERSION_PATCH := $(shell docker run --rm alpine/semver semver -c -i patch $(CURRENT_VERSION))
DOCKER_NEXT_VERSION_MINOR := $(shell docker run --rm alpine/semver semver -c -i minor $(CURRENT_VERSION))
DOCKER_NEXT_VERSION_MAJOR := $(shell docker run --rm alpine/semver semver -c -i major $(CURRENT_VERSION))

PROJECT = github.com/armory-io/arm
TAG=$(NEXT_VERSION)
BUILD_DIR=$(shell pwd)/build
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
#select all packages except a few folders because it's an integration test
PKGS := $(shell go list ./... | grep -v -e /integration -e /vendor)
CURRENT_DIR=$(shell pwd)
PROJECT_DIR_LINK=$(shell readlink ${PROJECT_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
# Go since 1.6 creates dynamically linked exes, here we force static and strip the result
# LINUX_LDFLAGS = -ldflags "-X ${PROJECT}/cmd.SEMVER=${TAG} -X ${PROJECT}/cmd.COMMIT=${COMMIT} -X ${PROJECT}/cmd.BRANCH=${BRANCH} -linkmode external -extldflags -static -s -w"
# LDFLAGS = -ldflags "-X ${PROJECT}/cmd.SEMVER=${TAG} -X ${PROJECT}/cmd.COMMIT=${COMMIT} -X ${PROJECT}/cmd.BRANCH=${BRANCH} -extldflags -s -w"

# Build the project
all: clean lint vet build

run:
	go run main.go

windows:
	GOOS=windows GOARCH=amd64 go build -o ${BUILD_DIR}/arm-${TAG}-windows-amd64.exe main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/arm-${TAG}-linux-amd64 main.go

darwin:
	GOOS=darwin GOARCH=amd64 go build -o ${BUILD_DIR}/arm-${TAG}-darwin-amd64 main.go

build: $(BUILD_DIR) windows linux darwin

test: 
	PCT=31 bin/test_coverage.sh

GOLINT=$(GOPATH)/bin/golint

$(GOLINT):
	go get -u golang.org/x/lint/golint

$(BUILD_DIR):
	mkdir -p $@
	chmod 777 $@

lint: $(GOLINT)
	@$(GOLINT) $(PKGS)

vet:
	go vet -v ./...

fmt:
	go fmt $$(go list ./... | grep -v /vendor/)

clean:
	rm -rf ${BUILD_DIR}
	go clean

.PHONY: lint test vet fmt clean run

current-version:
	@echo $(CURRENT_VERSION)

next-version-patch:
	@echo $(DOCKER_NEXT_VERSION_PATCH)

next-version-minor:
	@echo $(DOCKER_NEXT_VERSION_MINOR)

next-version-major:
	@echo $(DOCKER_NEXT_VERSION_MAJOR)

release-patch:
	git checkout master
	git tag $(DOCKER_NEXT_VERSION_PATCH)
	git push --tags --force

release-major:
	git checkout master
	git tag $(DOCKER_NEXT_VERSION_MINOR)
	git push --tags --force

release-minor:
	git checkout master
	git tag $(DOCKER_NEXT_VERSION_MAJOR)
	git push --tags --force
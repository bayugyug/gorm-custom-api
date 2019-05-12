BUILD_DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
BUILD_TIME := $(shell date +%Y%m%d.%H%M%S)
BUILD_HASH := $(shell git log -1 2>/dev/null| head -n 1 | cut -d ' ' -f 2)
TEST_FILES := $(shell go list ./... | grep -v /vendor/)
BUILD_COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_NAME := gorm-custom-api
BUILD_RELEASE?=0.0.1
PROJECT?=github.com/bayugyug/$(BUILD_NAME)

all: build

build : clean
	CGO_ENABLED=0 GOOS=linux go build -o bin/$(BUILD_NAME) -a -tags netgo -installsuffix netgo -installsuffix cgo -v -ldflags " \
	-X ${PROJECT}/configs.BuildTime=$(BUILD_TIME)  \
	-X ${PROJECT}/configs.Release=$(BUILD_RELEASE) \
	-X ${PROJECT}/configs.Commit=$(BUILD_COMMIT) " .

test : clean
	go test -v ./... > testrun.txt
	golint  $(TEST_FILES) > lint.txt
	go vet -v $(TEST_FILES) > vet.txt
	gocov test github.com/bayugyug/gorm-custom-api | gocov-xml > coverage.xml
	go test $(TEST_FILES) -bench=. -test.benchmem -v 2>/dev/null | gobench2plot > benchmarks.xml
	ginkgo -v  ./... > gink.txt

testginkgo : build
	ginkgo -v  ./...

testrun : clean test
	time go test -v -bench=. -benchmem -dummy > testrun.txt 2>&1

prepare : build

pack-alpine: clean build
	cd deployments && sudo docker build --no-cache --rm -t bayugyug/gorm-custom-api:alpine  -f  Dockerfile.alpine .

pack-scratch: clean build
	cp bin/$(BUILD_NAME) deployments/ && cd deployments && sudo docker build --no-cache --rm -t bayugyug/gorm-custom-api:scratch  -f  Dockerfile.scratch .

clean:
	rm -f $(BUILD_NAME) bin/$(BUILD_NAME) deployments/$(BUILD_NAME)
	rm -f benchmarks.xml coverage.xml vet.txt lint.txt testrun.txt gink.txt

re: clean all


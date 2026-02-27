export GOMOD=go.mod
export GORACE=halt_on_error=1
export GOFLAGS=
export CGO_ENABLED=1

DATE=$(shell date +'%Y.%m.%d.%H:%M:%S')
TARGET=tt

# for cloud build, if we do not have a build env variable set
# use the revision_id from the env
ifndef BUILD
override BUILD = 1.0
endif

LDFLAGS=-ldflags "-X main.Build=$(BUILD) -X main.BuildDate=$(DATE)"

.PHONY: all build build-race clean fmt vet bench test test-coverage lint

default: all

all: clean fmt lint vet test build

build:
	go build $(LDFLAGS) -o $(TARGET)

build-race:
	go build -race $(LDFLAGS) -o $(TARGET)

clean:
	rm -f $(TARGET)
	find . -type f -name "*.coverprofile" -delete
	find . -type f -name "*.out" -delete
	rm -rf vendor

fmt:
	go fmt $$(go list ./... | grep -v /vendor/)

vet:
	go vet $$(go list ./... | grep -v /vendor/)

bench:
	go test -bench=. -run=^$ -v ./...

test:
	go test -count=1 -race -v ./...

test-coverage:
	#lsof -i tcp:8000 | awk 'NR!=1 {print $2}' | xargs kill
	go test -count=1 -race -cover -coverprofile=coverage.out -coverpkg=./... -v ./... | grep '% of statements\|FAIL:' | grep 'ok\|FAIL:' | sed 's_github.com/regr76/__' | sed 's/ok\ \ \t/ /' | sed 's/\t[0-9]*.[0-9]*s\tcoverage:/ /' | sed 's/\ of\ statements/ /' > ./coverage/unit-test-coverage.txt
	go tool cover -func coverage.out | grep "total:" | sed 's/\t//g' |  sed 's/(statements)/ /' >> ./coverage/unit-test-coverage.txt
	# go tool cover -html="coverage.out"

lint:
	@golangci-lint run



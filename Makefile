include .env

PROJECTNAME=$(shell basename "$(PWD)")

# Go related variables.
GOBASE=$(shell pwd)
GOPATH="$(GOBASE)/vendor:$(GOBASE)"
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the gophermart
PID_GOPHER_MART=/tmp/.$(PROJECTNAME)-gopher-mart.pid

RANDOM=$(shell date +%s)
RND1=$(shell echo "("$RANDOM" % 1024) + 63490" | bc)
GOPHER_MART_PORT=$(RND1)
RND2=$(shell echo "("$RANDOM" % 1024) + 64514" | bc)
ACCRUAL_PORT=$(RND2)
ADDRESS=localhost:$(GOPHER_MART_PORT)
TEMP_FILE=$(shell mktemp)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## run: Compile and run server and agent
run: go-compile start

## start: Start in development mode. Auto-starts when code changes.
start: start-gopher-mart

## stop: Stop development mode. GOPHER_MART
stop: stop-gopher-mart

start-gopher-mart: stop-gopher-mart
	@echo "  >  $(PROJECTNAME) is available at $(ADDRESS)"
	@-$(GOBIN)/gophermart & echo $$! > $(PID_GOPHER_MART)
	@cat $(PID_GOPHER_MART) | sed "/^/s/^/  \>  PID: /"

stop-gopher-mart:
	@-touch $(PID_GOPHER_MART)
	@-kill `cat $(PID_GOPHER_MART)` 2> /dev/null || true
	@-rm $(PID_GOPHER_MART)

restart-gopher-mart: stop-gopher-mart start-gopher-mart

build: go-build-gophermart

## clean: Clean build files. Runs `go clean` internally.
clean:
	@(MAKEFILE) go-clean

## compile: Compile the binary.
go-compile: go-build-gopher-mart

go-build-gopher-mart:
	@echo "  >  Building gopher mart binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) cd ./cmd/gophermart && go build -o $(GOBIN)/gophermart $(GOFILES)

go-generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

.PHONY: go-update-deps
go-update-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

test:
	@echo "  > Test Iteration ..."
	go vet -vettool=$(which statictest) ./...
	cd bin && ./gophermarttest -test.v -test.run=^TestGophermart$$ -gophermart-binary-path=./gophermart -gophermart-host=localhost -gophermart-port=$(GOPHER_MART_PORT) -gophermart-database-uri="postgresql://postgres:postgres@localhost/praktikum?sslmode=disable" -accrual-binary-path=./accrual -accrual-host=localhost -accrual-port=$(ACCRUAL_PORT) -accrual-database-uri="postgresql://postgres:postgres@localhost/praktikum?sslmode=disable"

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

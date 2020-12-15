ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

# Database
POSTGRES_USER ?= user
POSTGRES_PASSWORD ?= password
POSTGRES_ADDRESS ?= localhost:5432
POSTGRES_DATABASE ?= test_user

# Migration Tools
MIGRATE_VERSION ?=v4.10.0
# Option:
# - darwin(Mac OS)
# - linux (choose this as the default since most of our server run on linux)
# - windows
MIGRATE_PLATFORM ?=linux
.PHONY: init
init: init-env migrate-prepare

.PHONY: init-env
init-env:

.PHONY: init-test
init-test: init
	@go get -v
	@go install -v github.com/jstemmer/go-junit-report

.PHONY: migrate-prepare
migrate-prepare:
	@rm -rf bin
	@mkdir bin

	# Reference: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#download-pre-built-binary-windows-macos-or-linux
	curl -L https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.$(MIGRATE_PLATFORM)-amd64.tar.gz | tar xvzO > ./bin/migrate
	chmod +x ./bin/migrate

.PHONY: migrate-create
migrate-create:
	@bin/migrate create -ext sql -dir repositories/migrations ${name}

.PHONY: test
test:
	@go test -v -race 1 ./...

.PHONY: e2e-test
e2e-test: init-test
	@go test -v -race ./e2e 2>&1 | $(GOPATH)/bin/go-junit-report > junit.xml

.PHONY: run
run:
	go run main.go

clean:
	@rm -rf user.db
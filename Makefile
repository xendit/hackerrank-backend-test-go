
# Database
POSTGRES_USER ?= user
POSTGRES_PASSWORD ?= password
POSTGRES_ADDRESS ?= localhost:5432
POSTGRES_DATABASE ?= test_user

# Migration Tools
MIGRATE_VERSION ?=v4.14.1
# Option:
# - darwin(Mac OS)
# - linux (choose this as the default since most of our server run on linux)
# - windows
MIGRATE_PLATFORM ?=darwin
.PHONY: init
init: init-env migrate-prepare
 
.PHONY: init-env
init-env:
	@sudo rm -rf /opt/go
	@wget -c https://golang.org/dl/go1.15.6.linux-amd64.tar.gz
	@sudo tar -C /opt/ -xzf go1.15.6.linux-amd64.tar.gz

.PHONY: init-test
init-test: init
	go get -u github.com/jstemmer/go-junit-report

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
	@go test -v -race -p 1 ./...

.PHONY: e2e-test
e2e-test:
	@go test -v -race -p 1  ./e2e 2>&1 | go-junit-report > junit.xml

.PHONY: run
run:
	go run main.go

clean:
	@rm -rf user.db
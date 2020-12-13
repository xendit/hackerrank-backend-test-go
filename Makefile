
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
	@curl -fsSL https://get.docker.com -o get-docker.sh
	@sudo sh get-docker.sh

	@sudo curl -L https://github.com/docker/compose/releases/download/1.18.0/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
	@sudo chmod +x /usr/local/bin/docker-compose
	@docker-compose --version

.PHONY: init-test
init-test: init
	@docker-compose up -d 

.PHONY: migrate-prepare
migrate-prepare:
	@rm -rf bin
	@mkdir bin
	# Reference: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#unversioned
	# @go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate
	# @go build -a -o ./bin/migrate -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate

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
	@go test -v -race -p 1  ./e2e

.PHONY: run
run:
	go run main.go

clean:
	@rm -rf user.db

# Database
POSTGRES_USER ?= user
POSTGRES_PASSWORD ?= password
POSTGRES_ADDRESS ?= localhost:5432
POSTGRES_DATABASE ?= test_user

# Migration Tools
MIGRATE_VERSION ?=v4.7.0
# Option:
# - darwin(Mac OS)
# - linux (choose this as the default since most of our server run on linux)
# - windows
MIGRATE_PLATFORM ?=darwin
.PHONY: init
init: migrate-prepare
 
.PHONY: migrate-prepare
migrate-prepare:
	# Reference: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#unversioned
	# @go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate
	# @go build -a -o ./bin/migrate -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate

	# Reference: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#download-pre-built-binary-windows-macos-or-linux
	curl -L https://github.com/golang-migrate/migrate/releases/download/$(MIGRATE_VERSION)/migrate.$(MIGRATE_PLATFORM)-amd64.tar.gz | tar xvzO > ./bin/migrate
	chmod +x ./bin/migrate

.PHONY: migrate-create
migrate-create:
	@bin/migrate create -ext sql -dir repositories/migrations ${name}

.PHONY: migrate-up
migrate-up:
	@bin/migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" \
	-path=repositories/migrations up

.PHONY: migrate-down
migrate-down:
	@bin/migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" \
	-path=repositories/migrations down

.PHONY: migrate-force
migrate-force:
	@bin/migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" \
	-path=repositories/migrations force ${version}

.PHONY: migrate-drop
migrate-drop:
	@bin/migrate -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_ADDRESS)/$(POSTGRES_DATABASE)?sslmode=disable" \
	-path=repositories/migrations drop

.PHONY: test
test:
	@go test -v -race -p 1 ./...

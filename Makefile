GREEN='\033[0;32m'
ORANGE='\033[0;33m'
NC='\033[0m'

PG_PORT=54320
test-db-up: ## Start postgres db in docker with PG_PORT
	@echo ${ORANGE}"Running the test DB"${NC}
	-docker stop gophermart-test-db
	docker run --rm -e POSTGRES_USER=gophermart -e POSTGRES_PASSWORD=1234 -p $(PG_PORT):5432 --name gophermart-test-db -d \
		postgres:12-alpine postgres -c log_statement=all
	docker exec gophermart-test-db timeout 20s bash -c "until pg_isready -d gophermart -U gophermart; do sleep 0.5; done"
	sleep 0.5 # need for mac
	# migrate -source file://migrations -database postgres://gophermart:1234@127.0.0.1:$(PG_PORT)/gophermart?sslmode=disable up
	@echo ${GREEN}"db up"${NC}

test-db-stop: ## Stop and clean postgres
	@echo ${ORANGE}"Stopping the test DB"${NC}
	# migrate -source file://migrations -database postgres://gophermart:1234@127.0.0.1:$(PG_PORT)/biller?sslmode=disable down -all
	docker stop gophermart-test-db
	@echo ${GREEN}"db down"${NC}


tidy: ## Format the code
	@echo ${ORANGE}"Formatting the code..."${NC}
	go mod tidy
	go fmt ./...
	@echo ${GREEN}"done"${NC}

run: ## Run the app
	ACCRUAL_SYSTEM_ADDRESS=localhost:8095
	DATABASE_URI=postgres://gophermart:1234@127.0.0.1:54320/gophermart?sslmode=disable
	RUN_ADDRESS=localhost:8094
	go run ./cmd/gophermart

build: ## Building the app
	go build ./cmd/gophermart

coverage: ## Checking test coverage
	go tool cover -func cover.out | grep total:

test: ## Launching tests
	go test ./...

clean: ## Cleaning building files
	rm gophermart

gen-mocks:
	cd internal && mockery --all --case underscore --keeptree --disable-version-string --with-expecter
GREEN='\033[0;32m'
ORANGE='\033[0;33m'
NC='\033[0m'

PG_PORT=54320
test-db-up: ## Start postgres db in docker with PG_PORT
	@echo ${ORANGE}"Running the test DB"${NC}
	-docker stop gophemart-test-db
	docker run --rm -e POSTGRES_USER=gophemart -e POSTGRES_PASSWORD=1234 -p $(PG_PORT):5432 --name gophemart-test-db -d \
		postgres:12-alpine postgres -c log_statement=all
	docker exec gophemart-test-db timeout 20s bash -c "until pg_isready -d gophemart -U gophemart; do sleep 0.5; done"
	sleep 0.5 # need for mac
	# migrate -source file://migrations -database postgres://gophemart:1234@127.0.0.1:$(PG_PORT)/gophemart?sslmode=disable up
	@echo ${GREEN}"db up"${NC}

test-db-stop: ## Stop and clean postgres
	@echo ${ORANGE}"Stopping the test DB"${NC}
	# migrate -source file://migrations -database postgres://gophemart:1234@127.0.0.1:$(PG_PORT)/biller?sslmode=disable down -all
	docker stop gophemart-test-db
	@echo ${GREEN}"db down"${NC}
migrate:
	@echo "running migrate up..."
	cd internal && cd sql && cd schema && goose postgres "postgres://postgres:password@localhost:5432/devpay?sslmode=disable" up

migratedown:
	@echo "running migrate down..."
	cd internal && cd sql && cd schema && goose postgres "postgres://postgres:password@localhost:5432/devpay?sslmode=disable" down

dbreset: migratedown migrate
sqlc:
	@echo "Generating sqlc queries..."
	sqlc generate

up_build:
	@echo "running project..."
	docker-compose down && docker-compose up --build -d

down:
	@echo "shutting down app..."
	docker-compose down

run:
	go run cmd/*.go

runtb:
	cd ../../.. && rm -rf 0_0.tigerbeetle && ./tigerbeetle format --cluster=0 --replica=0 --replica-count=1 0_0.tigerbeetle && ./tigerbeetle start --addresses=3000 0_0.tigerbeetle

resettb:
	cd ../../.. && rm -rf 0_0.tigerbeetle

run:
	go run ./cmd/api/main.go
swaggo:
	swag init -g **/**/*.go
test:
	go test -cover ./...
docker:
	docker build -t ms-pari-web .

docker-run:
	docker-compose up --build -d
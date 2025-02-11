run:
	docker-compose up --build

test:
	go test ./tests/unit/service/...
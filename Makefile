.PHONY: up down restart logs ps test

up:
	sudo docker-compose up --build -d

down:
	sudo docker-compose down

restart: down up

logs:
	sudo docker-compose logs -f

ps:
	sudo docker-compose ps

test:
	cd services/auth-service && go test -v ./...

build:
	go	mod vendor && cd deployments && docker-compose -p stonks build

run:
	cd deployments && docker-compose up --build

down:
	cd deployments && docker-compose down
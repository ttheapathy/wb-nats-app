rbuild:
	go mod vendor && docker-compose up --build

run:
	docker-compose up

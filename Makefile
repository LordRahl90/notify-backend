MAINFILE=main.go

start: build-image start-docker

.PHONY : start 

run:
	go run ${MAINFILE}

build-image:
	docker build -t lordrahl/notify-server:latest -f .docker/Dockerfile .

start-docker:
	docker-compose -f .docker/docker-compose.yml up

kill:
	docker-compose -f .docker/docker-compose.yml kill
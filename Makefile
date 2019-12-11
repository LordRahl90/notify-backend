MAINFILE=main.go

start: build-image start-docker
execute: validate deploy

.PHONY : start execute

run:
	go run ${MAINFILE}

build-image:
	docker build -t lordrahl/notify-server:latest -f .docker/Dockerfile .

push:
	docker push lordrahl/notify-server:latest

start-docker:
	docker-compose -f .docker/docker-compose.yml up

kill:
	docker-compose -f .docker/docker-compose.yml kill

validate:
	circleci config validate

deploy:
	circleci local execute --skip-checkout
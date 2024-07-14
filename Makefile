POSTGRES_SERVER_NAME=x_postgres
APP_SERVER_NAME=x_app
SERVER_NAME=x_backend

build:
	make down
	docker-compose build

up:
	docker-compose up -d

exec_app:
	docker-compose exec -it ${APP_SERVER_NAME} bash

exec_pg:
	docker-compose exec -it ${POSTGRES_SERVER_NAME} bash

down:
	docker-compose down --rmi all --volumes

stop:
	docker-compose stop

build_server:
	docker build -f ./docker/Dockerfile -t ${SERVER_NAME} .

run_server:
	docker run --name=${SERVER_NAME} ${SERVER_NAME}

stop_server:
	docker stop ${SERVER_NAME}

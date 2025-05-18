## INCLUDE ENV
#!make
include .env
export $(shell sed 's/=.*//' .env)

## VARS
dbs = "postgresql://$(LOCAL_DB_USER):$(LOCAL_DB_PWD)@$(LOCAL_DB_HOST):$(LOCAL_DB_PORT)/$(LOCAL_DB_DBNAME)?sslmode=disable"
migration_dir = "./db/migrations"


## CREATING AND MANAGING DB
dbinit:
	sudo docker run --name $(LOCAL_DB_DBNAME) -h localhost -p $(LOCAL_DB_PORT):$(LOCAL_DB_PORT) -e POSTGRES_USER=$(LOCAL_DB_USER) -e POSTGRES_PASSWORD=$(LOCAL_DB_PWD) postgres:latest
dbcreate:
	sudo docker exec -it $(LOCAL_DB_DBNAME) createdb --username=$(LOCAL_DB_USER) --owner=$(LOCAL_DB_USER) $(LOCAL_DB_DBNAME)

dbdrop:
	sudo docker exec -it $(LOCAL_DB_DBNAME) dropdb $(LOCAL_DB_DBNAME)


## To run DB after restart
# sudo docker ps -a | grep iztech_agsm
# sudo docker start <container_id>

## Migrations

gooseup:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dbs) goose -dir=$(migration_dir) up

goosedown:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dbs) goose -dir=$(migration_dir) down

goosereset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dbs) goose -dir=$(migration_dir) reset


## RUN APP
build:
	go mod tidy
	sudo docker build -f Dockerfile.auth -t iztech-agms-be-auth:latest .

run-auth:
	sudo docker run -p 8080:8080 iztech-agms-be-auth:latest


## UTILS
# goose create command ---> migrate -dir [DIR] create [MIGRATION_NAME] [DRIVER]
# goose create command ---> migrate -dir ./db/migrations create my_migration sql
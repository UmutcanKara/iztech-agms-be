## INCLUDE ENV
#!make
include .env
export $(shell sed 's/=.*//' .env)

# dbs = "postgresql://root:pwd@0.0.0.0:5432/chess?sslmode=disable"

## VARS
dbs = "postgresql://$(LOCAL_DB_USER):$(LOCAL_DB_PWD)@$(LOCAL_DB_HOST):$(LOCAL_DB_PORT)/$(LOCAL_DB_DBNAME)?sslmode=disable"
migration_dir = "./db/migrations"


## CREATING AND MANAGING DB
dbinit:
	sudo docker run --name $(LOCAL_DB_DBNAME) -p $(LOCAL_DB_PORT):$(LOCAL_DB_PORT) -e POSTGRES_USER=$(LOCAL_DB_USER) -e POSTGRES_PASSWORD=$(LOCAL_DB_PWD) postgres:latest
dbcreatedb:
	sudo docker exec -it $(LOCAL_DB_DBNAME) createdb --username=$(LOCAL_DB_USER) --owner=$(LOCAL_DB_USER) $(LOCAL_DB_DBNAME)

dbdropdb:
	sudo docker exec -it $(LOCAL_DB_DBNAME) dropdb $(LOCAL_DB_DBNAME)

gooseup:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dbs) goose -dir=$(migration_dir) up

goosedown:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dbs) goose -dir=$(migration_dir) down

goosereset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(dbs) goose -dir=$(migration_dir) reset


## RUN APP
build:
	go mod tidy
	sudo docker build -f Dockerfile.auth -t iztech-agms:latest .

run-auth:
	sudo docker run -p 8080:8080 iztech-agms:latest
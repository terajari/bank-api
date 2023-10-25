include .env

db_container=bank-api-db

postgres:
	docker run --name ${db_container} -p 5432:5432 -e POSTGRES_USER=${USERNAME} -e POSTGRES_PASSWORD=${PASSWORD} -e POSTGRES_DB=${DB_NAME} -d postgres

createdb:
	docker exec -it ${db_container} createdb --username=${USERNAME} --owner=${USERNAME} ${DB_NAME}

dropdb:
	docker exec -it ${db_container} dropdb --username=${USERNAME} ${DB_NAME}

createmigrate:
	migrate create -ext sql -dir migration/postgres init-schema  

migrateup:
	migrate -path migration/postgres/ -database "${DB_SOURCE}" -verbose up

migratedown:
	migrate -path migration/postgres/ -database "${DB_SOURCE}" -verbose down

.PHONY: postgres createdb dropdb createmigrate migrateup migratedown
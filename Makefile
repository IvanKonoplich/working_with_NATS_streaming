init_db:
	docker run --name nats_streaming_postgres -p 5436:5432 -e POSTGRES_PASSWORD=qwerty -e POSTGRES_USER=test -e POSTGRES_DB=test -d postgres

reset_db:
	docker stop nats_streaming_postgres
	docker rm nats_streaming_postgres



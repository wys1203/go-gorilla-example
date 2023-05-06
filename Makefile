.PHONY: run-postgres stop-postgres

run-postgres:
	docker run --name postgres-ui-test \
		-e POSTGRES_DB=ui_test \
		-e POSTGRES_USER=ui_test \
		-e POSTGRES_PASSWORD=mysecretpassword \
		-p 5432:5432 \
		-d postgres

stop-postgres:
	docker stop postgres-ui-test
	docker rm postgres-ui-test

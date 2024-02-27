build:
	go build -o build/sonify

clean:
	rm -f build/

run:
	air

migrate:
	goose -dir sql/schema postgres "user=postgres dbname=sonify sslmode=disable password=12345" up

migration-status:
	goose -dir sql/schema postgres "user=postgres dbname=sonify sslmode=disable password=12345" status

sqlgen:
	sqlc generate
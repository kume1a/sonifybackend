build:
	go build -o build/sonify

clean:
	rm -f build/

run:
	air

migrate:
	goose -dir sql/schema postgres "user=postgres dbname=sonify sslmode=disable password=123456" up

migration-status:
	goose -dir sql/schema postgres "user=postgres dbname=sonify sslmode=disable password=123456" status

sqlgen:
	sqlc generate
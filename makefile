# Run without building
run:
	go run .

# build binary
build:
	go build .

# run goose migration up
migrate-up:
	goose -dir sql/schema postgres "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable" up

# run goose migration down
migrate-down:
	goose -dir sql/schema postgres "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable" down
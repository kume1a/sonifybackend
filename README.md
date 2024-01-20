# Serve

install `air`
```bash
go install github.com/cosmtrek/air@latest
```

Run and watch file changes
```bash
air
```

# Migrations

```bash
goose postgres "user=postgres dbname=sonify sslmode=disable password=12345" status
goose postgres "user=postgres dbname=sonify sslmode=disable password=12345" up
```

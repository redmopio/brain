# brain

## Dependencies

- Go
- SQLC
- PostgreSQL

## Dev

```bash
# Download
go mod download

# Copy and fill env vars
cp .env.example .env

# Run
go run ./cmd/brain
```

## TODO

- [ ] Add timeouts and cancellations
- [ ] Create groups table

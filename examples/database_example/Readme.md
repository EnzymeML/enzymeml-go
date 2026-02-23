# v2 API Example

This example starts the EnzymeML v2 REST API using SQLite.

## Run

```bash
go run main.go
```

The server starts at `http://localhost:8080`.

Interactive docs and OpenAPI:

- `http://localhost:8080/v2/docs`
- `http://localhost:8080/v2/openapi.json`
- `http://localhost:8080/v2/openapi.yaml`

## Quick check

```bash
curl -X POST http://localhost:8080/v2/documents \
  -H "Content-Type: application/json" \
  -d '{"name":"my-doc","version":"2.0.0"}'

curl http://localhost:8080/v2/documents
curl http://localhost:8080/v2/openapi.json
```

The API includes CRUD endpoints for all v2 resources under `/v2`.

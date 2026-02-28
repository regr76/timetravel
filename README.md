## To Run The Server

1. Compile and run the Go application:
```bash
cd timetravel
make all
tt
make clean
```

2. Test the server using the healthcheck endpoint:
```bash
curl http://localhost:8000/health
```

You should see the following response:
```json
{"ok":true}
```


## The Assignment

A core part of any insurance platform is a reliable and auditable
record-keeping system. It must store all relevant data used to underwrite
policies. Policyholders periodically submit and update information about the
risks they want covered, such as their desired liability limits or changes to
their workforce. These changes can significantly affect the policy's risk
profile and, consequently, the premium.

The current codebase represents a very simplified version of this system, with:
- `GET /api/v1/record/{id}` – retrieves a record (a simple JSON mapping of
strings to strings)
- `POST /api/v1/record/{id}` – creates or updates a record

### Objective 1: Persist Data with SQLite

Replace the in-memory storage backend with a persistent SQLite database. The
goal is to ensure that all record data is retained even if the server is shut
down and restarted.

### Objective 2: Implement Time Travel Functionality

Introduce a “time travel” feature that allows querying the state of any record
at a specific point in time. This enables accurate reconstructions for
compliance, audits, and risk recalculations.

This objective is open-ended and may require significant changes across the
codebase. You'll introduce **record versioning and history tracking**.

Build out a new set of endpoints under `/api/v2` with the following
functionality:
- Retrieve records at specific versions (not just the latest)
- Apply updates to the latest version while preserving history
- List all available versions of a record
- Ensure full backward compatibility: `/api/v1` endpoints should continue to
work as-is, with no changes in behavior

## Reference -- V1 API

The current API consists of just two endpoints:
- `GET /api/v1/records/{id}`
- `POST /api/v1/records/{id}`,

All ids must be **positive integers**.

## Reference -- V2 API
- `GET /health`
- `GET /api/v1/records/{id}`
- `POST /api/v1/records/{id}`
- `GET /api/v1/records/{id}/versions/{version}`
- `GET /api/v1/records/{id}/list`
- TODO: `GET /api/v1/records/{id}/datetime`
,

All ids and versions must be **positive integers**.

### `GET /api/v1/records/{id}`

Retrieves a record by its ID. If the record exists, the server returns it in
JSON format. If the record does not exist, an error message is returned.

✅ Successful Response Example
```bash
> GET /api/v1/records/2323 HTTP/1.1

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8

{"id": 2323, "data": {"david": "hey", "davidx": "hey"}}
```

❌ Error Response Example
```bash
> GET /api/v1/records/32 HTTP/1.1

< HTTP/1.1 400 Bad Request
< Content-Type: application/json; charset=utf-8

{"error": "record of id 32 does not exist"}
```

### `POST /api/v1/records/{id}`

Creates or updates a record at the specified ID.
- If the record does not exist, it will be created.
- If the record already exists, it will be updated.
- Payload values must be a JSON object with string keys and values (or `null`).
- Keys with `null` values will be deleted from the record.

✅ Create a Record
```bash
> POST /api/v1/records/1 HTTP/1.1
> Content-Type: application/json

{"hello": "world"}

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8

{"id": 1, "data": {"hello": "world"}}
```

🔁 Update a Record
```bash
> POST /api/v1/records/1 HTTP/1.1
> Content-Type: application/json

{"hello": "world 2", "status": "ok"}

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8

{"id": 1, "data": {"hello": "world 2", "status": "ok"}}
```

❌ Delete a field from a record
```bash
> POST /api/v1/records/1 HTTP/1.1
> Content-Type: application/json

{"hello": null}

< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8

{"id": 1, "data": {"status": "ok"}}
```

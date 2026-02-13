# kvtxt

<p align="center">
  <b>Encrypted key-value storage over HTTP - single binary, zero framework</b>
</p>

<p align="center">

  <!-- Language -->
  <img src="https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go" />

  <!-- HTTP -->
  <img src="https://img.shields.io/badge/HTTP-net%2Fhttp-informational" />

  <!-- Storage -->
  <img src="https://img.shields.io/badge/Database-SQLite-003B57?logo=sqlite" />
  <img src="https://img.shields.io/badge/WAL-Enabled-success" />

  <!-- Encryption -->
  <img src="https://img.shields.io/badge/Encryption-AES--256--GCM-darkgreen" />

  <!-- Cache -->
  <img src="https://img.shields.io/badge/Cache-In--Memory%20LRU-orange" />
  <img src="https://img.shields.io/badge/TTL-Supported-success" />

  <!-- Build -->
  <img src="https://img.shields.io/badge/Build-Static%20Binary-blue" />

  <!-- License -->
  <img src="https://img.shields.io/badge/License-MIT-lightgrey" />

</p>

`kvtxt` is a minimal, API-only key–value storage service written in Go.

It accepts arbitrary UTF-8 payloads (including structured JSON), stores them **encrypted at rest (AES-256-GCM)**, and retrieves them using opaque, URL-safe keys.

Data is persisted in SQLite (WAL mode enabled). An optional in-memory cache reduces read latency. 

It is intentionally simple, explicit, and production-friendly.

---

## Design Goals

* Minimal surface area
* Clear separation of layers
* Deterministic behavior
* Encrypted storage by default
* Single static binary
* No hidden abstractions

---

## Architecture

Layered structure:

* `cmd/` → Application entry point
* `internal/api/` → HTTP layer (router, middleware, handlers)
* `internal/storage/` → SQLite persistence
* `internal/crypto/` → AES-256-GCM encryption
* `internal/cache/` → In-memory LRU cache
* `internal/config/` → Environment configuration
* `internal/worker/` → TTL cleanup worker

Request flow:

```
Client
   ↓
HTTP API
   ↓
Encrypt (AES-256-GCM)
   ↓
SQLite (WAL)
   ↕
In-memory Cache
```

---

## Features

* REST API (`POST`/`GET`)
* AES-256-GCM encryption at rest
* SQLite single file database
* WAL mode enabled 
* TTL/expiration
* TTL-aware LRU cache
* Opaque URL-safe keys
* Single static Go binary

---

## API

### Create Entry

**POST** `/v1/kv`

Request body:

```json
{
  "text": {
    "message": "hello world"
  },
  "content_type": "application/json",
  "ttl_seconds": 1000
}
```

Fields:

* `text` (required)
  Can be:

    * JSON object
    * JSON array
    * String
* `content_type` (optional)
* `ttl_seconds` (required)

Example:

```bash
curl --location 'http://localhost:8080/v1/kv' \
--header 'Content-Type: application/json' \
--data '{
    "text": {
        "message" : "hello world"
    },
    "content_type": "application/json",
    "ttl_seconds" : 1000
}'
```

Response:

```json
{
  "key": "adfXWRDY0TEFP6Zm",
  "expires_at": 1770915497
}
```

* `key` → Opaque identifier
* `expires_at` → Unix timestamp

---

### Retrieve Entry

**GET** `/v1/kv/{key}`

Example:

```bash
curl --location 'http://localhost:8080/v1/kv/adfXWRDY0TEFP6Zm'
```

Response:

```json
{
  "message": "hello world"
}
```

Behavior:

* If original payload was JSON - returned as JSON
* If original payload was plain text - returned as raw text
* Expired or unknown key - `404 Not Found`

---

## TTL Behavior

* `ttl_seconds` controls expiration.
* Expired entries:

    * Are treated as non-existent.
    * May be cleaned by background worker.
* If omitted, the entry does not expire.

---

## Configuration

Environment variables:

| Variable               | Description        | Default      |
| ---------------------- | ------------------ | ------------ |
| `KVTXT_PORT`           | HTTP bind address  | `:8080`      |
| `KVTXT_DB_PATH`        | SQLite file path   | `./kvtxt.db` |
| `KVTXT_ENCRYPTION_KEY` | 32-byte base64 key | required     |

Example:

```bash
export KVTXT_DB_PATH=./kvtxt.db
export KVTXT_ENCRYPTION_KEY=<base64-32-byte-secret>
```

Generate a secure key:

```bash
openssl rand -base64 32
```

---

## Build

### Build Binary

```bash
git clone <repo-url>
cd kvtxt

go mod tidy
go build -o kvtxt ./cmd/kvtxt
```

---

## Run

```bash
export KVTXT_DB_PATH=./kvtxt.db
export KVTXT_ENCRYPTION_KEY=<base64-32-byte-secret>

./kvtxt
```

Server starts at:

```
http://localhost:8080
```

---

## Security

* AES-256-GCM encryption at rest
* Encryption key never stored in DB
* Opaque keys reduce enumeration risk
* WAL improves durability
* No authentication (by design)

If exposed publicly, place behind:

* TLS termination
* Reverse proxy
* Authentication layer
* Rate limiting

---

## Performance

* Cache reduces read latency
* WAL enables concurrent reads
* Low memory footprint
* No reflection or heavy frameworks
* Suitable for small to medium workloads

---

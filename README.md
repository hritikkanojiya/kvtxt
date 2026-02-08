# kvtxt

`kvtxt` is a lightweight, API-only key–value text storage service exposed over HTTP.

It accepts arbitrary UTF-8 text payloads, stores them **encrypted at rest**, and retrieves them using opaque, URL-safe keys. Data is persisted in SQLite, with an in-memory cache used to reduce read latency.

There is **no UI**, **no authentication**, and **no external dependencies** beyond SQLite.
The service is designed to be boring, explicit, and production-friendly.

---

## Features

* Simple REST API (POST / GET)
* Encrypted payloads at rest (AES-256-GCM)
* SQLite persistent storage (single file)
* WAL mode enabled
* Optional TTL / expiry support
* In-process LRU cache with TTL awareness
* No frameworks, no ORM
* Single static Go binary
* Docker-ready (scratch / distroless target)
---

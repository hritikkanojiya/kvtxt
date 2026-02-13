# ---------- Builder ----------
FROM golang:1.24-bookworm AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y build-essential

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o kvtxt ./cmd/kvtxt


# ---------- Runtime ----------
FROM gcr.io/distroless/cc-debian12

WORKDIR /app

COPY --from=builder /app/kvtxt .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/kvtxt"]

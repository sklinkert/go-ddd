# syntax=docker/dockerfile:1

# ---- Build stage ----
FROM golang:1.26 AS build
WORKDIR /src

# Cache dependencies first.
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .

# Build statically-linked binaries (server + migration tool).
ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor
RUN go build -o /out/marketplace ./cmd/marketplace
RUN go build -o /out/migrate ./migrate.go

# ---- Runtime stage ----
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=build /out/marketplace /app/marketplace
COPY --from=build /out/migrate /app/migrate
# Migrations are read at runtime by the migrate tool.
COPY --from=build /src/migrations /app/migrations

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/marketplace"]

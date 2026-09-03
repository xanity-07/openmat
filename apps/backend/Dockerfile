# syntax=docker/dockerfile:1

FROM golang:1.27 AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

FROM base AS dev

RUN go install github.com/air-verse/air@latest

COPY . .

CMD ["air", "-c", ".air.toml"]

FROM base AS builder

COPY . .

RUN CGO_ENABLED=0 go build \
    -o /openmat \
    ./cmd/openmat/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS production

COPY --from=builder /openmat /openmat

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/openmat"]

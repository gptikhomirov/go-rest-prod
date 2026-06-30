FROM golang:1.26.3-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/go-rest-prod/exe /app/cmd/go-rest-prod

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /app/cmd/go-rest-prod/exe /app/exe
CMD ["/app/exe"]

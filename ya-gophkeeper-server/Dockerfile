FROM golang:1.22.2 AS builder
WORKDIR /build
COPY . .
RUN go test -cover ./...
RUN go build -o main cmd/app/main.go

FROM golang:1.22.2 AS ya-gophkeeper-server
RUN apt-get update
WORKDIR /build
COPY --from=builder /build/main /build/main
COPY --from=builder /build/db/migrations /build/migrations
COPY --from=builder /build/config.json /build/config.json
COPY --from=builder /build/tls /build/tls
EXPOSE 8080
EXPOSE 8081
CMD ["./main"]
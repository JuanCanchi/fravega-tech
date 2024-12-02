FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o /app/main main.go

FROM golang:1.20

WORKDIR /app

COPY --from=builder /app/main main

EXPOSE 8080

CMD ["/app/main"]

FROM golang:1.17 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.14 AS runner

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
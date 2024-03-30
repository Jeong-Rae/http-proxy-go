FROM golang:1.22.1 AS builder

WORKDIR /app

COPY go.mod ./
COPY proxy.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.14 AS runner

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
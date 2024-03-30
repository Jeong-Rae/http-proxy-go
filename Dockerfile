# 빌드 스테이지
FROM golang:1.17 AS builder

WORKDIR /app

# 소스 코드 복사
COPY proxy.go .

# 애플리케이션 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 런타임 스테이지
FROM alpine:3.14

WORKDIR /app

# 빌드 스테이지에서 빌드된 바이너리 복사
COPY --from=builder /app/main .

# 애플리케이션 실행
CMD ["./main"]
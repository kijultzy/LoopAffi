# ==========================================
# Stage 1: Builder
# ==========================================
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git (diperlukan beberapa modul Go)
RUN apk add --no-cache git

# Download dependensi lebih dulu (manfaatkan layer cache Docker)
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy seluruh source code
COPY . .

# Build binary — static linking untuk binary yang portable
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o main .

# ==========================================
# Stage 2: Runtime image minimalis
# ==========================================
FROM alpine:3.19

# ca-certificates untuk HTTPS ke Render Postgres
RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Jakarta

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]

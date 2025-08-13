# ---------- Build stage ----------
FROM golang:1.24-alpine AS build
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct

# Modules
COPY go.* ./
RUN go mod download

# Source
COPY . .

# Generate Swagger inside image (ensures docs exist)
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.3
RUN /go/bin/swag init -g cmd/main.go -o docs

# Sanity check
RUN ls -la && ls -la cmd

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /app/app ./cmd

# ---------- Runtime stage ----------
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /app/app /app/app
COPY --from=build /app/docs /app/docs
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/app"]

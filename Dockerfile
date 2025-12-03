# syntax=docker/dockerfile:1.6

ARG GO_VERSION=1.23.0

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

# 静态编译
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN go build -ldflags="-w -s" -o /build/wms-server cmd/server/main.go

FROM alpine:latest AS runtime

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/wms-server ./wms-server

# 设置时区
ENV TZ=Asia/Shanghai

EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

LABEL org.opencontainers.image.title="wms-server" \
      org.opencontainers.image.description="WMS Go backend server" \
      org.opencontainers.image.version="1.0.0" \
      maintainer="WMS Dev Team"

ENTRYPOINT ["./wms-server"]

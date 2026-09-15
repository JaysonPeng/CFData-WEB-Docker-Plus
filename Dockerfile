# ==========================================
# CFData-WEB Docker
# 正式构建 combined_refactor
# 支持 linux/amd64 + linux/arm64
# ==========================================

# ---------- Build ----------
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src/combined_refactor

RUN apk add --no-cache \
    ca-certificates \
    git

# 先复制正式版本的 Go 模块文件，利用 Docker 构建缓存
COPY combined_refactor/go.mod combined_refactor/go.sum ./

RUN go mod download

# 复制完整正式版本源码
COPY combined_refactor/ ./

# 构建 combined_refactor 正式 Web 版本
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /cfdata \
    .


# ---------- Runtime ----------
FROM alpine:3.22

# 程序放 /app，数据/配置放 /data，避免宿主机挂载覆盖程序
WORKDIR /data

RUN apk add --no-cache \
    ca-certificates \
    tzdata

ENV TZ=Asia/Shanghai
ENV CFDATA_CONFIG_PATH=/data/cfdata-config.json

COPY --from=builder /cfdata /app/cfdata

EXPOSE 13335

ENTRYPOINT ["/app/cfdata"]

CMD ["-port", "13335"]

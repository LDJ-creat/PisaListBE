FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 最终镜像
FROM alpine:latest

WORKDIR /app

# 从builder阶段复制编译好的程序
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

# 暴露端口
EXPOSE 8080

# 运行程序
CMD ["./main"] 
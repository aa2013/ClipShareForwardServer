FROM --platform=$BUILDPLATFORM debian:latest

WORKDIR /app

ARG TARGETARCH
ARG TARGETOS
ARG VERSION
COPY forward_server_${TARGETOS}_${TARGETARCH}_${VERSION} forward_server

RUN apt-get update && apt-get install -y tzdata && rm -rf /var/lib/apt/lists/* && chmod +x ./forward_server

# 暴露端口
EXPOSE 80 9283

# 启动命令
CMD ["./forward_server"]
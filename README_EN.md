# ClipShareForwardServer

This project is the relay server for [ClipShare](https://github.com/aa2013/ClipShare).

## docker-compose

For more details, see the usage guide:
[Forward Server Guide](https://clipshare.coclyun.top/usages/forward.html)

```yaml
version: "3"
services:
  clipshare-forward-server:
    image: coclyun/clipshare-forward-server:latest
    container_name: clipshare-forward-server
    restart: always
    ports:
      - "8180:80"
      - "9283:9283"
    environment:
      - TZ=Asia/Shanghai
    volumes:
      - ./data:/app/data
```

## Admin Web UI

The management page provides:

- real-time network speed, traffic, and connection monitoring
- `config.yaml` editing
- plan and key management in private mode
- runtime log viewing

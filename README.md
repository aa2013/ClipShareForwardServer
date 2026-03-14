该项目是 [ClipShare](https://github.com/aa2013/ClipShare) 的中转程序。
## docker-compose 部署
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
## 管理页面：
主要提供以下功能：

+ 简单的网络速度/流量/连接数的实时查看
+ `config.yaml` 配置文件修改
+ 私有模式下的`套餐` 和 `密钥` 的管理
+ 运行日志查看

![connections](docs/images/connection.png)
![login](docs/images/login.png)
![plan](docs/images/plan.png)
![settings](docs/images/settings.png)
![logs](docs/images/logs.png)

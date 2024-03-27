# coinswap-server

### Redis

安装/启动Redis

### 配置说明

在`./config/config.toml`文件中配置了默认的系统配置项，按需求修改。


### 构建docker镜像

```bash
docker build -t irisnet/coinswap-server .
```

### 启动镜像

```bash
docker run -d -v ./config:/root/.farm -p 8080:8080 irisnet/coinswap-server start
```
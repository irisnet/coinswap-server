# coinswap-server

# SetUp

# Build And Run

- Build: `make all`
- Run: `make run`
- Cross compilation: `make build-linux`

## Config Description
### [config.toml](https://github.com/irisnet/coinswap-server/blob/main/config/config.toml)

```text
[server]
address = ":8080"
price_denom = "htltbcbusd"
handle_farms_worker_num = 11

[mongodb]
node_uri= "mongodb://db_user:db_password@127.0.0.1:27017/?connect=direct&authSource=db_name"
database= "db_name"

[redis]
address = "127.0.0.1:6379"
db = 0
password = ""

[irishub]
chain_name ="Irishub"
chain_id = "test"
#optional
base_denom = "uiris"
fee = "100000uiris"
grpc_address = "192.168.150.60:29090"
rpc_address = "tcp://192.168.150.60:26657"
lcd_address = "http://192.168.150.60:1317"

[task]
enable = true
cron_time_update_total_volume_lock = 5
cron_time_update_liquidity_pool = 5
cron_time_update_farm = 5
```
## Run with docker
You can run application with docker.

### Image
- Build  Image

```$xslt
docker build -t irisnet/coinswap-server .
```

### Run Application

```bash
docker run -d -v ./config:/root/.farm -p 8080:8080 irisnet/coinswap-server start
```
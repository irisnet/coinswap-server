package monitor

import (
	"github.com/irisnet/coinswap-server/config"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"testing"
)

var (
	_meticsService MetricsService
)

func TestMain(m *testing.M) {
	//初始化redis连接
	redis.Connect("127.0.0.1:6379", "", 0)
	irishub.Connect("tcp://192.168.150.33:26657", "192.168.150.33:9090",
		"iris", "100000uiris")
	model.Init(&config.Config{
		MongoDb: config.Mongodb{
			NodeUri:  "mongodb://dapp:dapppassword@192.168.150.40:27017/?connect=direct&authSource=dapp-server",
			Database: "dapp-server",
		},
		Server: config.Server{
			PriceDenom: "ubusd",
		},
		IRIShub: config.IRIShub{
			LcdAddr: "http://192.168.150.33:1317",
		},
	})

	_meticsService = NewMetricsService()

	m.Run()
}
func TestMetricsService_SetCronTaskStatus(t *testing.T) {
	SetCronTaskStatus("hahha", 1)
	t.Log(_meticsService.QueryCronTaskStatus())
}

func TestMetricsService_QueryCronTaskStatus(t *testing.T) {
	t.Log(_meticsService.QueryCronTaskStatus())
}

func TestMetricsService_QueryLcdConnectionStatus(t *testing.T) {
	data := _meticsService.QueryLcdConnectionStatus()
	t.Log(data)
}

func TestMetricsService_QueryRedisConnectionStatus(t *testing.T) {
	data := _meticsService.QueryRedisConnectionStatus()
	t.Log(data)
}

func TestMetricsService_QueryIrishubConnectionStatus(t *testing.T) {
	data := _meticsService.QueryIrishubConnectionStatus()
	t.Log(data)
}

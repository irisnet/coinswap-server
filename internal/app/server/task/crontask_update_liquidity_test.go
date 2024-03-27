package task

import (
	"github.com/irisnet/coinswap-server/config"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	sdktypes "github.com/irisnet/irishub-sdk-go/types"
	"testing"
)

func TestMain(m *testing.M) {
	//初始化redis连接
	redis.Connect("127.0.0.1:6379", "", 0)

	//irishub.Connect("tcp://34.77.68.145:26657", "34.77.68.145:9090",
	//	"irishub-1", "100000uiris")
	//model.Init(&config.Config{
	//	MongoDb: config.Mongodb{
	//		NodeUri:  "mongodb://iris:irispassword@127.0.0.1:27018/?connect=direct&authSource=dapp-server",
	//		Database: "dapp-server",
	//	},
	//	Server: config.Server{
	//		PriceDenom: "htltbcbusd",
	//	},
	//})

	irishub.Connect("tcp://192.168.150.60:46657", "192.168.150.60:49090",
		"irishub-qa", "100000uiris")
	model.Init(&config.Config{
		MongoDb: config.Mongodb{
			NodeUri:  "mongodb://dapp:dappPassword@192.168.150.60:27017/?connect=direct&authSource=dapp-server",
			Database: "dapp-server",
		},
		Server: config.Server{
			PriceDenom: "htltbcbusd",
		},
	})

	m.Run()
}

func TestSyncLiquidityTask_UpdatePool(t *testing.T) {
	result, err := irishub.Coinswap().QueryAllPools(sdktypes.PageRequest{
		//this no use page and count
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err := new(SyncLiquidityTask).UpdatePool(result)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(types.MarshalJsonIgnoreErr(data))
}

func TestFarmPoolSyncer_getTokenPrice1(t *testing.T) {
	irisTokenAmt, _ := sdktypes.NewDecFromStr("607224.739905")
	otherTokenAmt, _ := sdktypes.NewDecFromStr("9884.761070")
	irisPrice, _ := sdktypes.NewDecFromStr("0.0000166667")
	tokenPrice := getTokenPrice1(irisTokenAmt, otherTokenAmt, irisPrice)
	t.Log("irisPrice:", irisPrice.String())
	t.Log("atomPrice:", tokenPrice.String())

	irisTokenAmt, _ = sdktypes.NewDecFromStr("6070.109034")
	otherTokenAmt, _ = sdktypes.NewDecFromStr("0.09983494")
	tokenPrice = getTokenPrice1(irisTokenAmt, otherTokenAmt, irisPrice)
	t.Log("bnbPrice:", tokenPrice.String())
}

func TestFarmPoolSyncer_getTokenPrice(t *testing.T) {
	irisTokenAmt, _ := sdktypes.NewDecFromStr("607224.739905")
	otherTokenAmt, _ := sdktypes.NewDecFromStr("9884.761070")
	irisTokenAmt1, _ := sdktypes.NewDecFromStr("6000.000001")
	busdTokenAmt, _ := sdktypes.NewDecFromStr("1")
	t.Log("irisPrice:", busdTokenAmt.Quo(irisTokenAmt1).String())
	tokenPrice := getTokenPrice2(irisTokenAmt, otherTokenAmt, irisTokenAmt1, busdTokenAmt)
	t.Log("atomPrice:", tokenPrice.String())

	irisTokenAmt, _ = sdktypes.NewDecFromStr("6070.109034")
	otherTokenAmt, _ = sdktypes.NewDecFromStr("0.09983494")
	tokenPrice = getTokenPrice2(irisTokenAmt, otherTokenAmt, irisTokenAmt1, busdTokenAmt)
	t.Log("bnbPrice:", tokenPrice.String())
}

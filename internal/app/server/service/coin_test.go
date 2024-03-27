package service

import (
	"fmt"
	"github.com/irisnet/coinswap-server/config"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	sdktype "github.com/irisnet/irishub-sdk-go/types"
	"testing"
)

func TestMain(m *testing.M) {
	//初始化redis连接
	redis.Connect("127.0.0.1:6379", "", 0)
	irishub.Connect("tcp://192.168.150.33:26657", "192.168.150.33:9090",
		"iris", "4uiris")
	model.Init(&config.Config{
		MongoDb: config.Mongodb{
			NodeUri:  "mongodb://iris:irispassword@127.0.0.1:27018/?connect=direct&authSource=dapp-server",
			Database: "dapp-server",
		},
		Server: config.Server{
			PriceDenom: "ubusd",
		},
		IRIShub: config.IRIShub{
			LcdAddr: "http://192.168.150.33:1317",
		},
	})

	m.Run()
}

func TestCoinService_QueryAllCoins(t *testing.T) {

	res, err := new(CoinService).QueryAllCoins()
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(res)

}

func TestCoinService_PriceByUSD(t *testing.T) {
	busdAmt := sdktype.OneDec()
	res, err := new(CoinService).PriceByUSD(busdAmt)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(res)
}

func TestCoinService_PriceIRISByUSD(t *testing.T) {
	uirisAmt := sdktype.OneDec()
	price, usdAmt, err := new(CoinService).PriceIRISByUSD(uirisAmt)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(price)
	fmt.Println(usdAmt)
}

func TestCoinService_QueryBaseCoinPrice(t *testing.T) {
	ret, err := new(CoinService).QueryBaseCoinPrice()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Pass")
	t.Log(ret)
}

func TestCoinService_GetUnitPrice(t *testing.T) {
	uniprice, _, err := new(CoinService).GetUnitPrice(sdktype.NewDecCoin("uiris", sdktype.NewInt(1)))
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("iris")
	t.Log(uniprice.String())
	bnbuniprice, _, err := new(CoinService).GetUnitPrice(sdktype.NewDecCoin("htltbcbnb", sdktype.NewInt(1)))
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("bnb")
	t.Log(bnbuniprice.String())
	busduniprice, _, err := new(CoinService).GetUnitPrice(sdktype.NewDecCoin("htltbcbusd", sdktype.NewInt(1)))
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("busd")
	t.Log(busduniprice.String())
	atomuniprice, _, err := new(CoinService).GetUnitPrice(sdktype.NewDecCoin("ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", sdktype.NewInt(1)))
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("atom")
	t.Log(atomuniprice.String())
}

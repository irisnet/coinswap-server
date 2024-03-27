package service

import (
	"encoding/json"
	"fmt"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/server/model"

	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"strconv"
	"testing"

	"github.com/irisnet/irishub-sdk-go/modules/coinswap"
	sdktype "github.com/irisnet/irishub-sdk-go/types"
)

func TestPoolService_CalPoolAPY(t *testing.T) {
	ps := PoolService{}
	p := &FarmPool{
		Staked:         3000000, //?
		RewardWeight:   1,
		Reward:         "atom",
		RewardPerBlock: 10,
	}
	f, t1, _ := ps.CalPoolAPY1(p, 1, 6307200)
	fmt.Println("apy:", f)
	fmt.Println(t1.String())
}

func (ps *PoolService) CalPoolAPY1(p *FarmPool,
	totalWeight int64, blockPerYear uint) (float64, sdktype.Dec, error) {
	lp := coinswap.QueryPoolResponse{
		Pool: sdktype.PoolInfo{
			Standard: sdktype.Coin{
				Denom:  "uiris",
				Amount: sdktype.NewInt(17811857513),
			},
			Token: sdktype.Coin{
				Denom:  "",
				Amount: sdktype.NewInt(1859688),
			},
			Lpt: sdktype.Coin{
				Denom:  "",
				Amount: sdktype.NewInt(808993733),
			},
			Fee: "",
		},
	}
	lpStaked := sdktype.NewDec(int64(p.Staked))
	uirisTotal := sdktype.NewDecFromInt(lp.Pool.Standard.Amount)
	tokenTotal := sdktype.NewDecFromInt(lp.Pool.Token.Amount)
	liquidityAmt := sdktype.NewDecFromInt(lp.Pool.Lpt.Amount)

	uirisStaked := lpStaked.Mul(uirisTotal).Quo(liquidityAmt)
	tokenStaked := lpStaked.Mul(tokenTotal).Quo(liquidityAmt)

	tokenStakedUiris := tokenStaked.Mul(uirisTotal).Quo(tokenTotal)
	totalUirisAmt := uirisStaked.Add(tokenStakedUiris)
	totalStaked := tokenStaked.Add(uirisStaked.Mul(tokenTotal).Quo(uirisTotal))

	totalIris := sdktype.NewDec(int64(p.RewardWeight)).
		QuoInt64(totalWeight).
		MulInt64(int64(p.RewardPerBlock)).
		MulInt64(int64(blockPerYear))

	//atomAmt := sdktype.NewDecFromInt(sdktype.NewInt(7560000000))
	aprDec := totalIris.Quo(totalStaked)
	fmt.Println("apr:", aprDec.String())

	var calAPY = func(apr sdktype.Dec) (apy float64, err error) {
		defer func() {
			if e := recover(); e != nil {
				apy = -1
				logger.Error("Redis Set failed", logger.Any("error", e))
				return
			}
		}()

		apyStr := apr.QuoInt64(365).Add(sdktype.OneDec()).Power(365).Sub(sdktype.OneDec()).String()
		apy, err = strconv.ParseFloat(apyStr, 64)
		if err != nil {
			return 0, err
		}
		return apy, nil
	}
	apy, err := calAPY(aprDec)
	//apy := math.Pow(1+tmp, 365) - 1
	return apy, totalUirisAmt, err
}

func TestPoolService_CheckClaim(t *testing.T) {
	has, err := new(PoolService).CheckClaim("iaa174qyl02cupyqq77cqqtdl0frda6dl3rp2h9snu")
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Pass")
	t.Log("isExist: ", has)
}

func TestPoolService_SaveClaim(t *testing.T) {
	if err := new(PoolService).SaveClaim("iaa174qyl02cupyqq77cqqtdl0frda6dl3rp2h9snu"); err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Pass")
}

func TestPoolService_QueryTotalVolumeLock(t *testing.T) {
	ret, err := new(PoolService).QueryTotalVolumeLock()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("Pass")
	t.Log("TVL:", ret)
}

func TestPoolService_QueryPoolInfo(t *testing.T) {

	//Convey("mock a farm ", t, func() {
	//
	//	//_, err := new(PoolService).QueryPoolInfo(model.Farms{}, "iaa174qyl02cupyqq77cqqtdl0frda6dl3rp2h9snu")
	//	//if err != nil {
	//	//	t.Fatal(err)
	//
	//	//}
	//
	//	resp, err := irishub.Coinswap().QueryAllPools(sdktype.PageRequest{})
	//	if err != nil {
	//		panic(err)
	//	}
	//	bytes, err := js.Marshal(resp.Pools)
	//	if err != nil {
	//		panic(err)
	//	}
	//	log.Println(string(bytes))
	//})

}

func TestPoolService_QueryAllPools(t *testing.T) {
	res, err := new(PoolService).QueryAllPools()
	if err != nil {
		t.Fatal(err.Error())
	}
	databytes, _ := json.Marshal(res)
	t.Log("Pass")
	t.Log(string(databytes))
}

func TestPoolService_QueryAllRewards(t *testing.T) {
	res, _, err := new(PoolService).QueryAllRewards("iaa174qyl02cupyqq77cqqtdl0frda6dl3rp2h9snu", "")
	if err != nil {
		t.Fatal(err.Error())
	}
	databytes, _ := json.Marshal(res)
	t.Log("Pass")
	t.Log(string(databytes))
}

func TestPoolService_CalPoolAPR(t *testing.T) {
	blocksPerYear, err := new(PoolService).BlocksPerYear()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("blocksPerYear:", blocksPerYear)
	farmPools, err := new(model.Farms).FindAllVisibleFarm()
	if err != nil {
		t.Fatal(err.Error())
	}
	pool, err := NewPoolFromFarm(&farmPools[0])
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err := new(PoolService).CalPoolAPR(farmPools[0], &pool, blocksPerYear)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("APR:", res.String())
	t.Log("VolumeLocked:", pool.VolumeLocked)
	//t.Log("FarmId:",farmPools[0].PoolId)
	//t.Log("data:",sdktype.OneDec().String())
	//value := res.QuoInt64(365).Add(sdktype.OneDec())
	//t.Log("value:",value.String())
	//data := sdktype.NewDec(value.TruncateInt64())
	//t.Log(data)
	//apyStr := res.QuoInt64(365).Add(sdktype.OneDec()).Power(365).Sub(sdktype.OneDec()).String()
	//t.Log("apyStr:", apyStr)
}

func TestPoolService_CalPoolAPY2(t *testing.T) {
	blocksPerYear, err := new(PoolService).BlocksPerYear()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("blocksPerYear:", blocksPerYear)
	farmPools, err := new(model.Farms).FindAllVisibleFarm()
	if err != nil {
		t.Fatal(err.Error())
	}
	var index int
	for i, val := range farmPools {
		if val.PoolId == "farm-35" {
			index = i
		}
	}
	pool, err := NewPoolFromFarm(&farmPools[index])
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err := new(PoolService).CalPoolAPY(farmPools[index], &pool, blocksPerYear)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("APY:", res)
}

func TestPoolService_FarmLptPrice(t *testing.T) {
	farmPools, err := new(model.Farms).FindAllVisibleFarm()
	if err != nil {
		t.Fatal(err.Error())
	}
	var index int
	for i, val := range farmPools {
		if val.PoolId == "farm-24" {
			index = i
		}
	}
	res, err := new(PoolService).FarmLptPrice(farmPools[index])
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("LptValue From APR:", res.String())

	farm := farmPools[index]
	liquidPoolResp, err := irishub.Coinswap().QueryPool(farm.TotalLptLocked.Denom)
	if err != nil {
		t.Fatal(err.Error())
	}
	farmLpAmout, err := sdktype.NewDecFromStr(farm.TotalLptLocked.Amount)
	if err != nil {
		t.Fatal(err.Error())
	}
	farmirisTotalAmt := farmLpAmout.Quo(sdktype.NewDecFromInt(liquidPoolResp.Pool.Lpt.Amount)).
		MulInt(liquidPoolResp.Pool.Standard.Amount)
	_, ValLock, err := new(PoolService).cs.PriceIRISByUSD(farmirisTotalAmt.MulInt64(2))
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("LptValue compute:", ValLock.String())
	t.Log(farm.PoolId, "diff", res.Sub(ValLock).String())
}

func TestPoolService_BlocksPerYear(t *testing.T) {
	blocksPerYear, err := new(PoolService).BlocksPerYear()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log("blocksPerYear:", blocksPerYear)
}

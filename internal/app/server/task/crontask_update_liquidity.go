package task

import (
	"fmt"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"github.com/irisnet/coinswap-server/internal/app/server/monitor"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	"github.com/irisnet/irishub-sdk-go/modules/coinswap"
	sdktypes "github.com/irisnet/irishub-sdk-go/types"
	"github.com/pkg/errors"
	"time"
)

func init() {
	RegisterTasks(&SyncLiquidityTask{})
}

type SyncLiquidityTask struct {
}

type liquidityTokenAmt struct {
	Token    sdktypes.Coin
	Standard sdktypes.Coin
}

func (task *SyncLiquidityTask) DoTask(fn func(string) chan bool) error {
	if !model.GetConf().Task.Enable {
		logger.Warn("CronTask Config Disable", logger.String("taskName", task.Name()))
		return nil
	}
	result, err := irishub.Coinswap().QueryAllPools(sdktypes.PageRequest{
		//this no use page and count
	})
	if err != nil {
		monitor.SetCronTaskStatus(types.SyncLiquidityTaskName, -1)
		return errors.Wrap(err, "Query all liquidity pool by sdk client failed")
	}

	liqutityPools, err := task.UpdatePool(result)
	if err != nil {
		monitor.SetCronTaskStatus(types.SyncLiquidityTaskName, -1)
		return errors.Wrap(err, "UpdatePool all liquidity pool failed")
	}
	if err := SaveOrUpdateStatisticInfo(types.MarshalJsonIgnoreErr(liqutityPools), types.StatisticLiquidityPool); err != nil {
		monitor.SetCronTaskStatus(types.SyncLiquidityTaskName, -1)
		return errors.Wrap(err, "save or update liquidity pool faild")
	}
	monitor.SetCronTaskStatus(types.SyncLiquidityTaskName, 1)
	return nil
}

func (task *SyncLiquidityTask) Name() string {
	return types.SyncLiquidityTaskName
}

func (task *SyncLiquidityTask) Cron() int {
	if model.GetConf().Task.CronTimeUpdateLiquidityPool > 0 {
		return model.GetConf().Task.CronTimeUpdateLiquidityPool
	}
	return types.CronTimeSyncLiquidityTask
}

func (task *SyncLiquidityTask) Start() {
	timeInterval := task.Cron()
	RunTimer(timeInterval, Sec, func() {
		monitor.SetCronTaskStatus(types.SyncLiquidityTaskName, 0)
		if err := task.DoTask(nil); err != nil {
			logger.Error(err.Error(), logger.String("taskName", task.Name()))
		}
	})

}

func getTokenPrice1(irisToken, otherToken, irisPrice sdktypes.Dec) sdktypes.Dec {
	return irisToken.Quo(otherToken).Mul(irisPrice)
}
func getTokenPrice2(irisToken, otherToken, irisToken1, usdToken sdktypes.Dec) sdktypes.Dec {
	return irisToken.Quo(otherToken).Mul(usdToken.Quo(irisToken1))
}

func getTokenPrice(irisToken, otherToken sdktypes.Coin, irisToken1, usdToken sdktypes.Int, coinScaleMap map[string]int) sdktypes.Dec {
	irisScale, okIrisScale := coinScaleMap[irisToken.Denom]
	otherScale, okOtherScale := coinScaleMap[otherToken.Denom]
	priceScale, okPriceScalue := coinScaleMap[model.GetConf().Server.PriceDenom]
	if okPriceScalue && okIrisScale && okOtherScale {
		return sdktypes.NewDecFromInt(irisToken.Amount).QuoInt(sdktypes.NewIntWithDecimal(1, irisScale)).
			Quo(sdktypes.NewDecFromInt(otherToken.Amount).QuoInt(sdktypes.NewIntWithDecimal(1, otherScale))).
			Mul(sdktypes.NewDecFromInt(usdToken).QuoInt(sdktypes.NewIntWithDecimal(1, priceScale)).
				Quo(sdktypes.NewDecFromInt(irisToken1).QuoInt(sdktypes.NewIntWithDecimal(1, irisScale))))
	}

	return sdktypes.ZeroDec()
}

func (task *SyncLiquidityTask) UpdatePool(result *coinswap.QueryAllPoolsResponse) ([]types.QueryLiquidityPoolResponse, error) {

	coinScaleSlice, err := new(model.Asset).QueryAllScale()
	if err != nil {
		return nil, err
	}
	scaleMapData := make(map[string]int, len(coinScaleSlice))
	for _, val := range coinScaleSlice {
		scaleMapData[val.Denom] = val.Scale
	}

	liqutityPoolSlice := make([]types.QueryLiquidityPoolResponse, 0, len(result.Pools))
	irisPrice := sdktypes.ZeroDec()
	usdToken := sdktypes.ZeroInt()
	irisToken := sdktypes.ZeroInt()
	tokenAmtSlice := make([]liquidityTokenAmt, 0, len(result.Pools))
	for _, p := range result.Pools {
		response := types.QueryLiquidityPoolResponse{
			Denom:     p.Token.Denom,
			Standard:  p.Standard,
			Token:     p.Token,
			Liquidity: p.Lpt,
			Fee:       p.Fee,
			UpdateAt:  time.Now().Unix(),
		}
		liqutityPoolSlice = append(liqutityPoolSlice, response)

		//compute iris price
		if standardScale, ok := scaleMapData[p.Standard.Denom]; ok {
			tokenScale, ok := scaleMapData[p.Token.Denom]
			if !ok {
				logger.Warn("skip compute liquidity pool price for cannot found token scale",
					logger.String("token", p.Token.Denom))
				continue
			}
			if p.Token.Denom == model.GetConf().Server.PriceDenom {
				usdToken = p.Token.Amount
				irisToken = p.Standard.Amount
				standardAmount := sdktypes.NewDec(p.Standard.Amount.Int64()).Quo(sdktypes.NewDecFromInt(sdktypes.NewIntWithDecimal(1, standardScale)))
				tokenAmount := sdktypes.NewDec(p.Token.Amount.Int64()).Quo(sdktypes.NewDecFromInt(sdktypes.NewIntWithDecimal(1, tokenScale)))
				if !standardAmount.IsZero() {
					irisPrice = tokenAmount.Quo(standardAmount)
					//logger.Debug("irisPrice = tokenAmount / irisAmount",
					//	logger.String("irisPrice", irisPrice.String()),
					//	logger.String("denom", p.Standard.Denom))
					key := fmt.Sprintf(types.RedisKeyPriceDenom, p.Standard.Denom)
					if err := redis.Set(key, irisPrice.String(), types.RedisKKeyPriceExpiration); err != nil {
						logger.Error("Save tokenPrice to redis failed",
							logger.String("token", p.Standard.Denom),
							logger.String("error", err.Error()))
					}
				}

			} else {
				tokenAmtSlice = append(tokenAmtSlice, liquidityTokenAmt{
					Token:    p.Token,
					Standard: p.Standard,
				})
			}
		}
	}

	//compute token price
	if !usdToken.IsZero() && !irisToken.IsZero() && len(tokenAmtSlice) > 0 {
		for _, val := range tokenAmtSlice {
			if !val.Token.Amount.IsZero() && !val.Standard.Amount.IsZero() {
				//tokenPrice := getTokenPrice1(val.Standard.Amount, val.Token.Amount, irisPrice)
				//tokenPrice := getTokenPrice(val.Standard.Amount, val.Token.Amount, irisToken, usdToken)
				tokenPrice := getTokenPrice(val.Standard, val.Token, irisToken, usdToken, scaleMapData)
				//logger.Debug("tokenPrice = (irisAmount / tokenAmount) * irisPrice",
				//	logger.String("tokenPrice", tokenPrice.String()),
				//	logger.String("denom", val.Token.Denom),
				//)
				key := fmt.Sprintf(types.RedisKeyPriceDenom, val.Token.Denom)
				if err := redis.Set(key, tokenPrice.String(), types.RedisKKeyPriceExpiration); err != nil {
					logger.Error("update token price  to redis failed",
						logger.String("token", val.Token.Denom),
						logger.String("error", err.Error()))
				}
			}
		}
	}

	return liqutityPoolSlice, nil
}

package service

import (
	"fmt"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	sdktype "github.com/irisnet/irishub-sdk-go/types"
	"github.com/pkg/errors"
)

type CoinService struct{}

func NewCoinService() CoinService {
	return CoinService{}
}
func (CoinService) DelCache(redisKey string) error {
	return redis.DeleteByKey(redisKey)
}
func (CoinService) QueryAllPrice() (prices types.Prices, err error) {
	coinScale, err := asset.QueryAllScale()
	if err != nil {
		return nil, err
	}
	for _, val := range coinScale {
		if val.Denom == model.GetConf().Server.PriceDenom {
			prices = append(prices, types.CoinPrice{
				Denom: val.Denom,
				Price: "1",
			})
			continue
		}
		key := fmt.Sprintf(types.RedisKeyPriceDenom, val.Denom)
		if price, err := redis.GetString(key); err == nil {
			item := types.CoinPrice{
				Denom: val.Denom,
				Price: price,
			}
			prices = append(prices, item)
		}
	}
	return prices, nil
}
func (CoinService) QueryAllCoins() (coins types.Coins, err error) {

	coinSlice, err := asset.FindAll()
	if err != nil {
		return nil, err
	}
	var baseChainName = types.DefaultIrishubChainName
	for _, val := range coinSlice {
		if val.Denom == model.GetConf().IRIShub.BaseDenom {
			baseChainName = val.Name
		}
	}

	loadIbcTransferInfos := func(values []model.IbcInfo, chainName string) []types.IbcInfoItem {
		var resData []types.IbcInfoItem
		for _, val := range values {
			item := types.IbcInfoItem{
				WithdrawInfo: types.InPath{
					Port:      val.InPath.Port,
					Channel:   val.InPath.Channel,
					ChainId:   model.GetConf().IRIShub.ChainID,
					ChainName: baseChainName,
				},
			}
			if len(val.Traces) > 0 {
				item.DepositInfo = append(item.DepositInfo, types.TraceInfo{
					Port:      val.Traces[0].Port,
					Channel:   val.Traces[0].Channel,
					ChainId:   val.Traces[0].ChainId,
					ChainName: chainName,
					Platform:  val.Traces[0].Platform,
					Denom:     val.Traces[0].Denom,
				})
			}
			resData = append(resData, item)
		}

		return resData
	}

	for _, coin := range coinSlice {
		coins = append(coins, types.CoinItem{
			Name:             coin.Name,
			Symbol:           coin.Symbol,
			Denom:            coin.Denom,
			Scale:            coin.Scale,
			CoinID:           coin.CoinID,
			Platform:         coin.Platform,
			Protocol:         coin.Protocol,
			LpToken:          coin.DenomLpt,
			Offline:          coin.Offline,
			Icon:             coin.Icon,
			Tips:             coin.Tips,
			IbcTransferInfos: loadIbcTransferInfos(coin.IbcTransferInfos, coin.Name),
		})
	}
	return coins, nil
}

func (cs CoinService) QueryBaseCoinPrice() (sdktype.Dec, error) {
	priceDenom := model.GetConf().Server.PriceDenom
	coin, err := asset.FindCoinByDenom(priceDenom)
	if err != nil {
		return sdktype.ZeroDec(), errors.Wrapf(err, "FindCoinByDenom(%v)", priceDenom)
	}
	result, err := irishub.Coinswap().QueryPool(coin.DenomLpt)
	if err != nil {
		return sdktype.ZeroDec(), errors.Wrapf(err, "query coinswap pool lptDenom(%v)", coin.DenomLpt)
	}

	uirisAmt := result.Pool.Standard.Amount
	usdAmt := result.Pool.Token.Amount

	if !uirisAmt.GT(sdktype.ZeroInt()) {
		return sdktype.ZeroDec(), fmt.Errorf("the amount of %s is less than 0 in pool: %s",
			result.Pool.Standard.Denom, model.GetConf().Server.PriceDenom)
	}
	return sdktype.NewDecFromInt(usdAmt).QuoInt(uirisAmt), nil
}

func (cs CoinService) PriceIRISByUSD(uirisAmt sdktype.Dec) (sdktype.Dec, sdktype.Dec, error) {
	irisPrice, err := cs.QueryBaseCoinPrice()
	if err != nil {
		return sdktype.ZeroDec(), sdktype.ZeroDec(), err
	}

	//logger.Info("Compute current iris price", logger.String("price", irisPrice.String()))
	usdAmt, err := cs.PriceByUSD(irisPrice.Mul(uirisAmt))
	if err != nil {
		return sdktype.ZeroDec(), sdktype.ZeroDec(), err
	}
	return irisPrice, usdAmt, nil
}

func (CoinService) PriceByUSD(usdAmt sdktype.Dec) (sdktype.Dec, error) {
	priceDenom := model.GetConf().Server.PriceDenom

	var scale int64
	coin, err := asset.FindCoinByDenom(priceDenom)

	if err != nil {
		return sdktype.ZeroDec(), err
	}

	scale = int64(coin.Scale)

	//logger.Info("Compute usd price: " + usdAmt.String())
	precision := sdktype.NewDecFromIntWithPrec(sdktype.OneInt(), scale)
	usd := usdAmt.Mul(precision)
	//logger.Info("Compute usd price: " + usd.String())
	return usd, nil
}

//func (cs CoinService) GetPrice(coin sdktype.DecCoin) (sdktype.Dec, error) {
//	priceDenom := model.GetConf().Server.PriceDenom
//	res, err := asset.FindCoinByDenom(priceDenom)
//	if err != nil {
//		return sdktype.ZeroDec(), err
//	}
//	srcLp, err := irishub.Coinswap().QueryPool(res.DenomLpt)
//	if err != nil {
//		return sdktype.ZeroDec(), err
//	}
//	//logger.Info("src lp Pool information",
//	//	logger.String("BaseCoin", srcLp.Pool.Standard.String()),
//	//	logger.String("TokenCoin", srcLp.Pool.Token.String()),
//	//	logger.String("Liquidity", srcLp.Pool.Lpt.String()),
//	//)
//	resCoin, err := asset.FindCoinByDenom(coin.Denom)
//	if err != nil {
//		return sdktype.ZeroDec(), err
//	}
//	dstLp, err := irishub.Coinswap().QueryPool(resCoin.DenomLpt)
//	if err != nil {
//		return sdktype.ZeroDec(), err
//	}
//
//	srcRate := srcLp.Pool.Token.Amount.Mul(sdktype.NewIntWithDecimal(1, 18)).Quo(srcLp.Pool.Standard.Amount)
//	dstRate := dstLp.Pool.Token.Amount.Mul(sdktype.NewIntWithDecimal(1, 18)).Quo(dstLp.Pool.Standard.Amount)
//	srcDecRate := sdktype.NewDecFromInt(srcRate)
//	dstDecRate := sdktype.NewDecFromInt(dstRate)
//	price := coin.Amount.Mul(srcDecRate.Quo(dstDecRate))
//	formula := fmt.Sprintf("%s * (%s/%s)=%s * (%s/%s)=%s",
//		priceDenom,
//		srcLp.Pool.Token.Denom,
//		dstLp.Pool.Token.Denom,
//		coin.Amount.String(),
//		srcRate.String(),
//		dstRate.String(),
//		price.String(),
//	)
//	logger.Info("Calculate relative prices " + "formula" + formula)
//	return price, nil
//}

func (cs CoinService) GetPriceV2(coin sdktype.DecCoin) (sdktype.Dec, error) {

	logger.Debug("ready to calculate unit price", logger.String("denom", coin.Denom))
	unitPrice, scale, err := cs.GetUnitPrice(coin)
	if err != nil {
		return sdktype.ZeroDec(), err
	}

	logger.Debug("unit price calculation finished", logger.String("denom", coin.Denom), logger.String("unit price", unitPrice.String()), logger.String("scale", scale.String()))

	if scale.LTE(sdktype.ZeroDec()) {
		return sdktype.Dec{}, errors.New("scale can not be less than zero")
	}

	logger.Info("coinPrice", logger.String("unitPrice", unitPrice.String()), logger.Any("coin", coin))
	return coin.Amount.Quo(scale).Mul(unitPrice), nil
}

func (cs CoinService) GetUnitPrice(coin sdktype.DecCoin) (sdktype.Dec, sdktype.Dec, error) {

	switch coin.Denom {

	case model.GetConf().IRIShub.BaseDenom:
		return cs.GetStandardTokenUnitPrice(coin)

	case model.GetConf().Server.PriceDenom:
		return cs.GetPriceTokenUnitPrice(coin)

	default:
		return cs.GetTokenUnitPrice(coin)
	}
}

func (cs CoinService) GetPriceTokenUnitPrice(coin sdktype.DecCoin) (sdktype.Dec, sdktype.Dec, error) {

	priceCoin, err := asset.FindCoinByDenom(coin.Denom)
	if err != nil {
		return sdktype.ZeroDec(), sdktype.ZeroDec(), errors.Wrapf(err, "Asset.FindCoinByDenom(%v)", coin.Denom)
	}

	priceCoinScale := sdktype.NewDecFromInt(sdktype.NewIntWithDecimal(1, priceCoin.Scale))

	return sdktype.OneDec(), priceCoinScale, nil
}

func (cs CoinService) GetTokenUnitPrice(decCoin sdktype.DecCoin) (price sdktype.Dec, scale sdktype.Dec, err error) {

	priceDenom := model.GetConf().Server.PriceDenom

	priceCoin, err := asset.FindCoinByDenom(priceDenom)
	if err != nil {
		return price, scale, errors.Wrapf(err, "Asset.FindCoinByDenom(%v)", priceDenom)
	}
	priceLiquidPool, err := irishub.Coinswap().QueryPool(priceCoin.DenomLpt)
	if err != nil {
		return price, scale, errors.Wrapf(err, "CoinSwap.QueryPool(%v)", priceCoin.DenomLpt)
	}

	standardCoin, err := asset.FindCoinByDenom(model.GetConf().IRIShub.BaseDenom)
	if err != nil {
		return price, scale, errors.Wrapf(err, "Asset.FindCoinByDenom(%v)", model.GetConf().IRIShub.BaseDenom)
	}

	normalCoin, err := asset.FindCoinByDenom(decCoin.Denom)
	if err != nil {
		return price, scale, errors.Wrapf(err, "Asset.FindCoinByDenom(%v)", decCoin.Denom)
	}

	normalLiquidPool, err := irishub.Coinswap().QueryPool(normalCoin.DenomLpt)
	if err != nil {
		return price, scale, errors.Wrapf(err, "CoinSwap.QueryPool(%v)", normalCoin.DenomLpt)
	}

	var (
		// scale
		priceTokenScale = sdktype.NewDecFromInt(sdktype.NewIntWithDecimal(1, priceCoin.Scale))

		standardTokenScale = sdktype.NewDecFromInt(sdktype.NewIntWithDecimal(1, standardCoin.Scale))

		normalTokenScale = sdktype.NewDecFromInt(sdktype.NewIntWithDecimal(1, normalCoin.Scale))

		// amount
		priceTokenAmountInPoolU = sdktype.NewDecFromInt(priceLiquidPool.Pool.Token.Amount).Quo(priceTokenScale)

		standardTokenAmountInPoolU = sdktype.NewDecFromInt(priceLiquidPool.Pool.Standard.Amount).Quo(standardTokenScale)

		normalTokenAmountInPoolT = sdktype.NewDecFromInt(normalLiquidPool.Pool.Token.Amount).Quo(normalTokenScale)

		standardTokenAmountInPoolT = sdktype.NewDecFromInt(normalLiquidPool.Pool.Standard.Amount).Quo(standardTokenScale)
	)

	if standardTokenAmountInPoolU.LTE(sdktype.ZeroDec()) || normalTokenAmountInPoolT.LTE(sdktype.ZeroDec()) {
		logger.Info("denominator cannot be less than zero", logger.Any("decCoin", decCoin))
		return price, scale, nil
	}

	// ( IRISTokenAmount /  CTokenAmount ) *(BUSDTokenAmount / IRISTokenAmount )
	return standardTokenAmountInPoolT.Quo(normalTokenAmountInPoolT).Mul(priceTokenAmountInPoolU.Quo(standardTokenAmountInPoolU)), normalTokenScale, nil
}

func (cs CoinService) GetStandardTokenUnitPrice(decCoin sdktype.DecCoin) (price sdktype.Dec, scale sdktype.Dec, err error) {

	priceDenom := model.GetConf().Server.PriceDenom

	priceCoin, err := asset.FindCoinByDenom(priceDenom)
	if err != nil {
		return price, scale, errors.Wrapf(err, "Assert.FindCoinByDenom(%v)", priceDenom)
	}

	standardCoin, err := asset.FindCoinByDenom(decCoin.Denom)
	if err != nil {
		return price, scale, errors.Wrapf(err, "Assert.FindCoinByDenom(%v)", priceDenom)
	}

	liquidPoolResp, err := irishub.Coinswap().QueryPool(priceCoin.DenomLpt)
	if err != nil {
		return price, scale, errors.Wrapf(err, "CoinSwap.QueryPool(%v)", priceCoin.DenomLpt)
	}

	var (
		priceTokenScale = sdktype.NewDecFromInt(sdktype.NewIntWithDecimal(1, priceCoin.Scale))

		standardTokenScale = sdktype.NewDecFromInt(sdktype.NewIntWithDecimal(1, standardCoin.Scale))

		priceTokenAmountInPoolU = sdktype.NewDecFromInt(liquidPoolResp.Pool.Token.Amount).Quo(priceTokenScale)

		standardTokenAmountInPoolU = sdktype.NewDecFromInt(liquidPoolResp.Pool.Standard.Amount).Quo(standardTokenScale)
	)

	if standardTokenAmountInPoolU.LTE(sdktype.ZeroDec()) {

		logger.Info("denominator cannot be less than zero", logger.Any("decCoin", decCoin))
		return price, scale, nil
	}
	return priceTokenAmountInPoolU.Quo(standardTokenAmountInPoolU), standardTokenScale, nil
}

package service

import (
	"context"
	"fmt"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub/farm"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/pkg/errors"
	"strconv"
	"sync"
	"time"

	sdktype "github.com/irisnet/irishub-sdk-go/types"

	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
)

type PoolService struct {
	cs CoinService
}

func NewPoolService(cs CoinService) PoolService {
	return PoolService{
		cs: cs,
	}
}

func (ps *PoolService) QueryAllPools() ([]types.Pool, error) {
	// basic pool
	pools, err := ps.QueryPoolBasicInfo()
	if err != nil {
		return nil, err
	}
	return pools, nil
}

func (ps *PoolService) QueryAllRewards(address string, poolId string) ([]types.PoolReward, string, error) {
	var resp []types.PoolReward
	// fill pool with account
	rewards, err := ps.QueryFarmerInfo(address)
	if err != nil {
		logger.Error("get farmer's  reward occurs err  ",
			logger.Any("err", err),
			logger.String("address", address))
		// query rewards failed not return error for account which not stake lp tokens
		return resp, "", nil
	}

	whitelist, err := whiteList.FindAll()
	if err != nil {
		return nil, "", err
	}
	whiteMap := make(map[string]string, len(whitelist))
	for _, val := range whitelist {
		whiteMap[val.PoolId] = val.PoolName
	}

	for _, poolAccountInfo := range rewards.List {

		if _, ok := whiteMap[poolAccountInfo.PoolId]; !ok {
			//skip for only display whiteList farm pool
			continue
		}

		//fillter by poolId
		if len(poolId) > 0 && poolAccountInfo.PoolId != poolId {
			continue
		}

		var earned = make([]types.CoinStr, 0)
		for _, reward := range poolAccountInfo.PendingReward {
			earned = append(earned, types.CoinStr{
				Denom:  reward.Denom,
				Amount: reward.Amount.String(),
			})
		}

		resp = append(resp, types.PoolReward{
			Staked: poolAccountInfo.Locked.Amount.Uint64(),
			Earned: earned,
			PoolId: poolAccountInfo.PoolId,
		})

	}

	return resp, rewards.Height, nil
}

func (ps *PoolService) QueryPoolBasicInfo() (pools []types.Pool, err error) {

	if err := redis.GetObject(types.RedisKeyFarmPool, &pools); err == nil {
		return pools, nil
	}

	farmPools, err := new(model.Farms).FindAllVisibleFarm()
	if err != nil {
		return nil, fmt.Errorf("FindAllVisibleFarm occurs err:%v", err.Error())
	}

	if len(farmPools) == 0 {
		return pools, nil
	}

	blocksPerYear, latestHeight, err := ps.BlocksPerYear()
	if err != nil {
		return nil, fmt.Errorf("BlocksPerYear occurs err:%v", err.Error())
	}
	pools, err = ps.FarmsToPools(farmPools, blocksPerYear, latestHeight)
	if err != nil {
		return pools, fmt.Errorf("FarmsToPools occurs err:%v", err.Error())
	}

	if err := redis.SetObject(types.RedisKeyFarmPool, pools, types.RedisKeyFarmExpiration); err != nil {
		logger.Error("set redis key occurs err", logger.String("key", types.RedisKeyFarmPool), logger.Any("val", pools))
	}
	return pools, nil
}

func (ps *PoolService) FarmsToPools(farms []model.Farms, blocksPerYear int64, latestHeight int64) (pools []types.Pool, err error) {

	lpTokens, err := lpt.FindAll()
	if err != nil {
		return nil, errors.Wrap(err, "get all coin icon_lpt failed")
	}
	lpIconMap := make(map[string]string, len(lpTokens))
	for _, val := range lpTokens {
		lpIconMap[val.Denom] = val.Icon
	}
	defaultWorkNum := int(11)
	if model.GetConf().Server.HandleFarmsWorkerNum > 0 {
		defaultWorkNum = model.GetConf().Server.HandleFarmsWorkerNum
	}
	poolMap, err := DoHandleFarms(defaultWorkNum, farms, NewPoolFromFarm)
	if err != nil {
		return nil, errors.Wrap(err, "DoHandleFarms failed")
	}

	for _, farm := range farms {
		pool := poolMap[farm.PoolId]
		if icon, ok := lpIconMap[farm.TotalLptLocked.Denom]; ok {
			pool.Icon = icon
		}

		if farm.Expired {
			pool.Status = 1
			pool.ActivityOn = false
			pools = append(pools, pool)
			continue
		}

		if startHeight, err := strconv.ParseInt(farm.StartHeight, 10, 64); err == nil {
			pool.ActivityOn = startHeight <= latestHeight
		}

		//logger.Debug("ready to calculate apy", logger.String("poo_id", farm.PoolId), logger.Int64("blockPerYear", blocksPerYear))
		pool.APY, err = ps.CalPoolAPY(farm, &pool, blocksPerYear)
		if err != nil {
			return nil, fmt.Errorf("CalPoolAPY occurs err:%w", err)
		}
		logger.Debug("apy calculation finished", logger.String("pool_id", farm.PoolId), logger.Float64("apy", pool.APY))
		pools = append(pools, pool)
	}
	return
}

// NewPoolFromFarm  create a  pool with basic attributes  from a farm
func NewPoolFromFarm(farm *model.Farms) (types.Pool, error) {

	liquidPool, err := irishub.Coinswap().QueryPool(farm.TotalLptLocked.Denom)
	if err != nil {
		return types.Pool{}, errors.Wrapf(err, "CoinSwap.QueryPool(%v)", farm.TotalLptLocked.Denom)
	}

	pool := types.Pool{
		ID:   farm.PoolId,
		Name: farm.Name,
		Code: liquidPool.Pool.Token.Denom,
	}

	// reward
	var sep = ""
	for _, reword := range farm.RewardPerBlock {
		pool.Reward += sep + reword.Denom
		sep = ","
	}
	return pool, nil
}

func DoHandleFarms(workNum int, farms []model.Farms, dowork func(farm *model.Farms) (types.Pool, error)) (map[string]types.Pool, error) {
	retChan := make(chan types.Pool, len(farms))
	errChan := make(chan error, len(farms))
	var wg sync.WaitGroup
	wg.Add(workNum)
	for i := 0; i < workNum; i++ {
		num := i
		go func(num int) {
			defer func() {
				wg.Done()
			}()
			for id := range farms {
				if id%workNum != num {
					continue
				}
				pool, err := dowork(&farms[id])
				retChan <- pool
				errChan <- err
			}
		}(num)
	}
	wg.Wait()
	close(errChan)
	close(retChan)

	var err error
	for e := range errChan {
		if e != nil && err == nil {
			err = e
		}
	}
	if err != nil {
		return nil, err
	}

	poolMap := make(map[string]types.Pool, len(farms))
	for pool := range retChan {
		poolMap[pool.ID] = pool
	}
	return poolMap, nil
}

func (ps *PoolService) CalPoolAPY(farm model.Farms, pool *types.Pool, blocksPerYear int64) (float64, error) {

	//if farm.IsLocalFarm {
	//	return 0, nil
	//}

	//logger.Debug("ready to calculate apr", logger.String("pool_id", farm.PoolId))
	apr, err := ps.CalPoolAPR(farm, pool, blocksPerYear)
	if err != nil {
		return 0, fmt.Errorf("CalPoolAPR occurs err:%w pool_id:%s", err, farm.PoolId)
	}
	logger.Debug("apr calculation finished", logger.String("pool_id", farm.PoolId), logger.String("apr val", apr.String()))

	farmLptAmount, err := sdktype.NewDecFromStr(farm.TotalLptLocked.Amount)
	if err != nil {
		return 0, fmt.Errorf("farm.TotalLptLocked.Amount NewDecFromStr occurs err:%w pool_id:%s", err, farm.PoolId)
	}
	if farmLptAmount.IsZero() {
		// APY should return -1 when farm no lock lptoken
		return -1, nil
	}
	var calAPY = func(apr sdktype.Dec) (apy float64, err error) {
		defer func() {
			if e := recover(); e != nil {
				apy = -1
				logger.Error("CalPoolAPY failed", logger.Any("error", e),
					logger.String("pool_id", farm.PoolId))
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
	//	(1+APR/365)^365 - 1
	apy, err := calAPY(apr)
	return apy, err
}

// CalPoolAPR #Reference http://wiki.bianjie.ai/display/IRISHub/Farm
func (ps *PoolService) CalPoolAPR(farm model.Farms, pool *types.Pool, blockPerYear int64) (sdktype.Dec, error) {

	//logger.Debug("ready to calculate totalRewardToken price", logger.String("pool_id", farm.PoolId))
	totalRewardTokenPrice, err := ps.FarmTotalRewardTokenPrice(farm, sdktype.NewDec(blockPerYear))
	if err != nil {
		return sdktype.ZeroDec(), errors.Wrap(err, "FarmTotalRewardTokenPrice")
	}

	logger.Debug(" totalReardToken price calculation finished", logger.String("pool_id", farm.PoolId),
		logger.String("totalRewardTokenPrice", totalRewardTokenPrice.String()))

	//logger.Debug("ready to calculate farm lptTotalValue", logger.String("pool_id", farm.PoolId))
	lptPrice, err := ps.FarmLptPrice(farm)
	if err != nil {
		return sdktype.ZeroDec(), errors.Wrap(err, "FarmLptPrice")
	}

	logger.Debug(" lptTotalValue calculation finished ", logger.String("pool_id", farm.PoolId), logger.String("lptTotalValue", lptPrice.String()))

	lptPriceFloatVal, err := strconv.ParseFloat(lptPrice.String(), 64)
	if err != nil {
		return sdktype.ZeroDec(), errors.Wrapf(err, "strconv.ParserFloat(%v,64)", lptPrice.String())
	}
	pool.VolumeLocked = lptPriceFloatVal

	if lptPrice.IsZero() {
		return sdktype.ZeroDec(), nil
	}

	return totalRewardTokenPrice.Quo(lptPrice), nil
}

// FarmTotalRewardTokenPrice #Reference http://wiki.bianjie.ai/pages/viewpage.action?pageId=58060394
// numerator of arp formula
func (ps *PoolService) FarmTotalRewardTokenPrice(farm model.Farms, blocks sdktype.Dec) (sdktype.Dec, error) {

	// total token price
	var tokenEarnedPerYear = sdktype.ZeroDec()

	for _, reward := range farm.RewardPerBlock {

		rewardAmount, err := sdktype.NewDecFromStr(reward.Amount)
		if err != nil {
			return tokenEarnedPerYear, errors.Wrapf(err, "sdktype.NewDecFromStr(%v)", rewardAmount)
		}

		logger.Debug(" ready to calculate reward price", logger.String("pool_id", farm.PoolId), logger.Any("denom", reward.Denom), logger.String("amount", reward.Amount))
		rewardPrice, err := ps.cs.GetPriceV2(sdktype.NewDecCoinFromDec(reward.Denom, rewardAmount))

		if err != nil {
			return tokenEarnedPerYear, err
		}

		logger.Debug(" reward price calculation finished", logger.String("pool_id", farm.PoolId), logger.String("price", rewardPrice.String()), logger.Any("denom", reward.Denom), logger.String("amount", reward.Amount))

		tokenEarnedPerYear = tokenEarnedPerYear.Add(rewardPrice.Mul(blocks))
	}

	return tokenEarnedPerYear, nil

}

// FarmLptPrice #Reference http://wiki.bianjie.ai/pages/viewpage.action?pageId=58060394
// denominator  of arp formula
func (ps *PoolService) FarmLptPrice(farm model.Farms) (sdktype.Dec, error) {

	liquidPoolResp, err := irishub.Coinswap().QueryPool(farm.TotalLptLocked.Denom)
	if err != nil {
		return sdktype.ZeroDec(), errors.Wrapf(err, "irishub.CoinSwap.QueryPool(%s)", farm.TotalLptLocked.Denom)
	}

	farmLptAmount, err := sdktype.NewDecFromStr(farm.TotalLptLocked.Amount)
	if err != nil {
		return sdktype.ZeroDec(), errors.Wrapf(err, "sdktype.NewDecFromStr(%v)", farm.TotalLptLocked.Amount)
	}

	var liquidPool = liquidPoolResp.Pool
	logger.Debug(" ready to lpt pool standard token  price", logger.String("pool_id", farm.PoolId), logger.String("liquid_pool", farm.TotalLptLocked.Denom), logger.Any("denom", liquidPool.Standard.Denom), logger.String("amount", liquidPool.Standard.Amount.String()))

	tokenPrice, err := ps.cs.GetPriceV2(sdktype.NewDecCoin(liquidPool.Standard.Denom, liquidPool.Standard.Amount))
	if err != nil {
		return sdktype.ZeroDec(), err
	}
	logger.Debug("   lpt pool standard token price calculation finished", logger.String("price", tokenPrice.String()), logger.String("pool_id", farm.PoolId), logger.String("liquid_pool", farm.TotalLptLocked.Denom), logger.Any("denom", liquidPool.Standard.Denom), logger.String("amount", liquidPool.Standard.Amount.String()))
	// (lpt/totalLpt)*TokenAPrice *2

	logger.Info("calculate lpt total value", logger.String("lpt amount", farmLptAmount.String()),
		logger.String("pool Amount", liquidPoolResp.Pool.Lpt.Amount.String()),
		logger.String("token price", tokenPrice.String()),
	)
	return farmLptAmount.QuoInt(liquidPoolResp.Pool.Lpt.Amount).Mul(tokenPrice).Mul(sdktype.NewDec(2)), nil
}

func (ps *PoolService) BlocksPerYear() (int64, int64, error) {
	block, err := irishub.Block(context.Background(), nil)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "irishub.Client.Block(%v)", nil)
	}
	height100 := block.Block.Height - 100
	block100, err := irishub.Block(context.Background(), &height100)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "irishub.Client.Block(%v)", height100)
	}
	var avgTime = sdktype.ZeroDec()
	diffHeight := block.Block.Height - block100.Block.Height
	if diffHeight >= 100 {
		avgTime = sdktype.NewDec(block.Block.Time.Unix() - block100.Block.Time.Unix()).QuoInt64(100)
	} else {
		avgTime = sdktype.NewDec(block.Block.Time.Unix() - block100.Block.Time.Unix()).QuoInt64(diffHeight)
	}

	blocksPerYear := sdktype.NewDec(365 * 24 * 3600).QuoInt64(avgTime.TruncateInt64())

	return blocksPerYear.TruncateInt64(), block.Block.Height, nil
}

//Deprecated
//func (ps *PoolService) CalPoolAPYV1(p FarmPool, blockPerYear uint) (float64, sdktype.Dec, error) {
//	coin, err := new(model.Asset).FindCoinByDenom(p.Code)
//	if err != nil {
//		//logger.Error("Find Coin By Denom failed", "error", err.Error())
//		return 0, sdktype.ZeroDec(), err
//	}
//	lp, err := irishub.Coinswap().QueryPool(coin.DenomLpt)
//	if err != nil {
//		//logger.Error("Query liquidity pool failed", "error", err.Error())
//		return 0, sdktype.ZeroDec(), err
//	}
//
//	logger.Info("Lp Pool information",
//		logger.String("BaseCoin", lp.Pool.Standard.Amount.String()),
//		logger.String("TokenCoin", lp.Pool.Token.Amount.String()),
//		logger.String("Liquidity", lp.Pool.Lpt.Amount.String()),
//	)
//
//	logger.Info("Lp staked information",
//		logger.Uint64("StakedLp", p.Staked),
//	)
//	lpStaked := sdktype.NewDec(int64(p.Staked))
//	uirisTotal := sdktype.NewDecFromInt(lp.Pool.Standard.Amount)
//	tokenTotal := sdktype.NewDecFromInt(lp.Pool.Token.Amount)
//	liquidityAmt := sdktype.NewDecFromInt(lp.Pool.Lpt.Amount)
//
//	uirisStaked := lpStaked.Mul(uirisTotal).Quo(liquidityAmt)
//	tokenStaked := lpStaked.Mul(tokenTotal).Quo(liquidityAmt)
//
//	totalIrisAmtStaked := uirisStaked.Add(tokenStaked.Mul(uirisTotal).Quo(tokenTotal))
//	totalTokenAmtStaked := tokenStaked.Add(uirisStaked.Mul(tokenTotal).Quo(uirisTotal))
//
//	//totalRewardPerYear := sdktype.NewDec(int64(p.RewardWeight)).
//	//	QuoInt64(totalWeight).
//	//	MulInt64(int64(p.RewardPerBlock)).
//	//	MulInt64(int64(blockPerYear))
//
//	var aprDec sdktype.Dec
//	switch {
//	case p.Reward == p.Code:
//		aprDec = totalRewardPerYear.Quo(totalTokenAmtStaked)
//	case p.Reward == model.GetConf().IRIShub.BaseDenom:
//		aprDec = totalRewardPerYear.Quo(totalIrisAmtStaked)
//	default:
//		total, err := ps.cs.GetPrice(sdktype.NewDecCoin(p.Reward, totalRewardPerYear.RoundInt()))
//		if err != nil {
//			logger.Error("Get price failed", logger.String("reward", p.Reward),
//				logger.String("error", err.Error()))
//			return 0, sdktype.ZeroDec(), err
//		}
//
//		staked, err := ps.cs.GetPrice(sdktype.NewDecCoin(p.Code, totalTokenAmtStaked.RoundInt()))
//		if err != nil {
//			logger.Error("Get price failed", logger.String("reward", p.Reward),
//				logger.String("error", err.Error()))
//			return 0, sdktype.ZeroDec(), err
//		}
//		aprDec = total.Quo(staked)
//	}
//	logger.Info("APR information ", logger.String("apr", aprDec.String()))
//	var calAPY = func(apr sdktype.Dec) (apy float64, err error) {
//		defer func() {
//			if e := recover(); e != nil {
//				apy = -1
//				logger.Error("Redis Set failed", logger.String("error", err.Error()))
//				return
//			}
//		}()
//
//		apyStr := apr.QuoInt64(365).Add(sdktype.OneDec()).Power(365).Sub(sdktype.OneDec()).String()
//		apy, err = strconv.ParseFloat(apyStr, 64)
//		if err != nil {
//			return 0, err
//		}
//		return apy, nil
//	}
//	apy, err := calAPY(aprDec)
//	//apy := math.Pow(1+tmp, 365) - 1
//	return apy, totalIrisAmtStaked, err
//}

func (ps *PoolService) QueryTotalVolumeLock() (string, error) {
	return statisticModel.FindTotalVolumeLock()
}

func (ps *PoolService) SaveClaim(address string) error {
	return claimAddressModel.Save(model.ClaimAddress{
		Address:    address,
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
	})
}

func (ps *PoolService) CheckClaim(address string) (bool, error) {
	return claimAddressModel.Exist(address)
}

// QueryFarmerInfo query farmer's accountInfo  from chain ,with poolId and address
func (ps *PoolService) QueryFarmerInfo(address string) (*farm.QueryFarmerResponse, error) {
	var req = farm.QueryFarmerRequest{
		Farmer: address,
	}

	farmerResponse, err := irishub.Farm().Farmer(context.Background(), req)
	if err != nil {
		return nil, errors.Wrapf(err, "Farm.Farmer(%+v)", req)
	}
	return farmerResponse, nil
}

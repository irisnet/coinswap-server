package cache

import (
	"fmt"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	"strings"
	"time"
)

const (
	RedisKeyTotalVolumeLock           = "volume_total_locked"
	RedisKeyLiquidityDenom            = "liquidity:%s"
	RedisKeyTotalVolumeLockExpiration = 5 * time.Second
	RedisKeyLiquidityDenomExpiration  = 5 * time.Second
)

var _ model.IStatisticInfo = StatisticInfoCache{}

type StatisticInfoCache struct {
	statisticInfo model.StatisticInfo
}

func (s StatisticInfoCache) FindTotalVolumeLock() (string, error) {
	volume, err := redis.GetString(RedisKeyTotalVolumeLock)
	if err == nil {
		return volume, nil
	}
	volume, err = s.statisticInfo.FindTotalVolumeLock()
	if err != nil {
		return "", err
	}
	if err := redis.Set(RedisKeyTotalVolumeLock, volume, RedisKeyTotalVolumeLockExpiration); err != nil {
		logger.Error("Redis Set failed", logger.String("error", err.Error()))
	}
	return volume, nil
}

func (s StatisticInfoCache) FindLiquidityPoolByDenom(denom string) (types.QueryLiquidityPoolResponse, error) {
	var response types.QueryLiquidityPoolResponse
	key := fmt.Sprintf(RedisKeyLiquidityDenom, denom)
	err := redis.GetObject(key, &response)
	if err == nil {
		return response, nil
	}
	response, err = s.statisticInfo.FindLiquidityPoolByDenom(denom)
	if err != nil {
		return response, err
	}
	if err := redis.SetObject(key, response, RedisKeyLiquidityDenomExpiration); err != nil {
		logger.Error("Save liquidity pool to redis failed",
			logger.String("denom", denom),
			logger.String("error", err.Error()))
	}

	return response, nil
}

func (s StatisticInfoCache) FindLiquidityPoolByDenoms(denoms []string) ([]types.QueryLiquidityPoolResponse, error) {
	var response []types.QueryLiquidityPoolResponse
	key := fmt.Sprintf(RedisKeyLiquidityDenom, strings.Join(denoms, ","))
	err := redis.GetObject(key, &response)
	if err == nil {
		return response, nil
	}
	response, err = s.statisticInfo.FindLiquidityPoolByDenoms(denoms)
	if err != nil {
		return response, err
	}
	if err := redis.SetObject(key, response, RedisKeyLiquidityDenomExpiration); err != nil {
		logger.Error("Save liquidity pool to redis failed",
			logger.String("denom", strings.Join(denoms, ",")),
			logger.String("error", err.Error()))
	}
	return response, nil
}

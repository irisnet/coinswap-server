package cache

import (
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
)

const (
	RedisKeyLpTokens           = "lp_tokens"
	RedisKeyLpTokensExpiration = 0
)

var _ model.ILpTokens = LpTokensCache{}

type LpTokensCache struct {
	lpt model.LpTokens
}

func (l LpTokensCache) FindAll() ([]model.LpTokens, error) {
	var coins []model.LpTokens
	err := redis.GetObject(RedisKeyLpTokens, &coins)
	if err == nil {
		return coins, nil
	}
	coins, err = l.lpt.FindAll()
	if err != nil {
		return nil, err
	}
	if err := redis.SetObject(RedisKeyLpTokens, coins, RedisKeyLpTokensExpiration); err != nil {
		logger.Error("Save white list to redis failed",
			logger.String("error", err.Error()))
	}
	return coins, nil
}

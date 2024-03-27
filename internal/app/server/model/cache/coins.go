package cache

import (
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"github.com/qiniu/qmgo"
)

const (
	RedisKeyCoins           = "coins"
	RedisKeyCoinsExpiration = 0
)

var _ model.IAsset = AssetCache{}

type AssetCache struct {
	asset model.Asset
}

func (a AssetCache) QueryAllScale() ([]model.Asset, error) {
	return a.findAll()
}
func (a AssetCache) FindAll() ([]model.Asset, error) {
	return a.findAll()
}
func (a AssetCache) findAll() ([]model.Asset, error) {
	var coins []model.Asset
	err := redis.GetObject(RedisKeyCoins, &coins)
	if err == nil {
		return coins, nil
	}
	coins, err = a.asset.FindAll()
	if err != nil {
		return nil, err
	}
	if err := redis.SetObject(RedisKeyCoins, coins, RedisKeyCoinsExpiration); err != nil {
		logger.Error("Save white list to redis failed",
			logger.String("error", err.Error()))
	}

	return coins, nil
}

func (a AssetCache) FindCoinByDenom(denom string) (model.Asset, error) {
	coins, err := a.findAll()
	if err != nil {
		return model.Asset{}, err
	}
	for _, val := range coins {
		if val.Denom == denom {
			return val, nil
		}
	}
	return model.Asset{}, qmgo.ErrNoSuchDocuments
}

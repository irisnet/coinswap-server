package cache

import (
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
)

const (
	RedisKeyWhiteList           = "farm:white_list"
	RedisKeyWhiteListExpiration = 0
)

var _ model.IWhiteList = WhiteListCache{}

type WhiteListCache struct {
	whiteList model.WhiteList
}

func (w WhiteListCache) FindAll() ([]model.WhiteList, error) {
	var whiteList []model.WhiteList
	err := redis.GetObject(RedisKeyWhiteList, &whiteList)
	if err == nil {
		return whiteList, nil
	}
	whiteList, err = w.whiteList.FindAll()
	if err != nil {
		return nil, err
	}
	if err := redis.SetObject(RedisKeyWhiteList, whiteList, RedisKeyWhiteListExpiration); err != nil {
		logger.Error("Save white list to redis failed",
			logger.String("error", err.Error()))
	}

	return whiteList, nil
}

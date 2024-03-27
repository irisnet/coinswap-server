package cache

import (
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
)

const (
	RedisKeyClaim = "claims"
)

var _ model.IClaimAddress = ClaimAddressCache{}

type ClaimAddressCache struct {
	claimAddress model.ClaimAddress
}

func (c ClaimAddressCache) Exist(address string) (bool, error) {
	if exist, err := redis.HExists(RedisKeyClaim, address); err == nil && exist {
		return true, nil
	}
	return c.claimAddress.Exist(address)
}

func (c ClaimAddressCache) Save(info model.ClaimAddress) error {
	if exist, err := redis.HExists(RedisKeyClaim, info.Address); err == nil && exist {
		return nil
	}
	err := c.claimAddress.Save(info)
	if err != nil {
		return err
	}

	if err := redis.HSet(RedisKeyClaim,
		info.Address, true); err != nil {
		logger.Error("Redis Set failed", logger.String("error", err.Error()))
	}
	return nil
}

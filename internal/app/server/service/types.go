package service

import (
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"github.com/irisnet/coinswap-server/internal/app/server/model/cache"
)

type FarmPool struct {
	ID             string `boil:"id" json:"id"`
	Code           string `boil:"code" json:"code"`
	Name           string `boil:"name" json:"name"`
	Staked         uint64 `boil:"staked" json:"staked"`
	RewardWeight   uint   `boil:"reward_weight" json:"reward_weight"`
	Reward         string `boil:"reward" json:"reward"`
	RewardPerBlock uint64 `boil:"reward_per_block" json:"reward_per_block"`
	Status         int    `boil:"status" json:"status"`
}

var (
	statisticModel    model.IStatisticInfo = new(cache.StatisticInfoCache)
	claimAddressModel model.IClaimAddress  = new(cache.ClaimAddressCache)
	asset             model.IAsset         = new(cache.AssetCache)
	lpt               model.ILpTokens      = new(cache.LpTokensCache)
	whiteList         model.IWhiteList     = new(cache.WhiteListCache)
	notice            model.Notice
)

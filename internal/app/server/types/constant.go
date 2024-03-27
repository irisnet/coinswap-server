package types

import "time"

const (
	SyncLiquidityTaskName             = "SyncLiquidityTask"
	SyncUpdateTotalVolumnLockTaskName = "SyncUpdateTotalVolumnLockTask"
	SyncFarmTaskName                  = "SyncFarmTaskName"
	CronTimeSyncLiquidityTask         = 2
	CronTimeUpdateTotalVolumnLock     = 5
	CronTimeSyncFarmTask              = 5

	StatisticTVL            = "TVL"
	StatisticLiquidityPool  = "Liquidity"
	DefaultIrishubChainName = "Irishub"

	RedisKeyPriceDenom       = "coins:price:%s"
	RedisKKeyPriceExpiration = 5 * time.Second

	RedisKeyFarmPool       = "farm"
	RedisKeyFarmExpiration = 5 * time.Second

	RedisKeyLatestBlock           = "latest_block"
	RedisKeyLatestBlockExpiration = 5 * time.Second
)

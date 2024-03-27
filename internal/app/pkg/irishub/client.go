package irishub

import (
	"context"
	"fmt"
	"github.com/irisnet/coinswap-server/config"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub/farm"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	sdk "github.com/irisnet/irishub-sdk-go"
	"github.com/irisnet/irishub-sdk-go/modules/bank"
	"github.com/irisnet/irishub-sdk-go/modules/coinswap"
	"github.com/irisnet/irishub-sdk-go/modules/keys"
	sdktype "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/irishub-sdk-go/types/store"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

var (
	client sdk.IRISHUBClient

	maxLimit = 100
)

// NewWalletService return a instance of  WalletService
func Connect(rcpAddr, grpcAddr, chainID, fee string) {
	options := []sdktype.Option{
		sdktype.KeyDAOOption(store.NewMemory(nil)),
		sdktype.TimeoutOption(10),
		sdktype.CachedOption(true),
	}
	if len(fee) > 0 {
		feeCoins, err := sdktype.ParseDecCoins(fee)
		if err != nil {
			panic(err)
		}
		options = append(options, sdktype.FeeOption(feeCoins))
	}

	cfg, err := sdktype.NewClientConfig(rcpAddr, grpcAddr, chainID, options...)
	if err != nil {
		panic(err)
	}

	client = sdk.NewIRISHUBClient(cfg)

}

func GetClient() sdk.IRISHUBClient {
	return client
}
func Block(ctx context.Context, height *int64) (*ctypes.ResultBlock, error) {
	var res ctypes.ResultBlock
	key := types.RedisKeyLatestBlock
	if height != nil {
		key = fmt.Sprintf(types.RedisKeyLatestBlock+":%d", *height)
	}
	if err := redis.GetObject(key, &res); err == nil {
		return &res, nil
	}
	ret, err := client.Block(ctx, height)
	if err != nil {
		return nil, err
	}
	if err := redis.SetObject(key, ret, types.RedisKeyLatestBlockExpiration); err != nil {
		logger.Error("set redis key occurs err", logger.String("key", key), logger.Any("val", ret))
	}
	return ret, nil
}
func Keys() keys.Client {
	return client.Key
}

func Coinswap() coinswap.Client {
	return client.Swap
}

func Bank() bank.Client {
	return client.Bank
}

func Farm() farm.Client {
	//return farm.NewLcdFarmClient(config.Get().IRIShub.LcdAddr)
	return farm.NewLcdFarmCacheClient(config.Get().IRIShub.LcdAddr)
}

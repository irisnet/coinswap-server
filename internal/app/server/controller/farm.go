package controller

import (
	"context"
	"github.com/irisnet/coinswap-server/internal/app/pkg/kit"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/server/service"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	sdktype "github.com/irisnet/irishub-sdk-go/types"
	"github.com/pkg/errors"
)

type farmController struct {
	BaseController
	service.CoinService
	service.PoolService
}

func (fc farmController) GetEndpoints() (endpoints []kit.Endpoint) {
	endpoints = append(endpoints, kit.Endpoint{
		URI:     "/farm/coins",
		Method:  "GET",
		Handler: fc.makeHandler(fc.QueryCoinList, nil),
	})

	endpoints = append(endpoints, kit.Endpoint{
		URI:     "/farm/pools",
		Method:  "GET",
		Handler: fc.makeHandler(fc.QueryFarmPoolList, nil),
	})

	endpoints = append(endpoints, kit.Endpoint{
		URI:     "/farm/rewards",
		Method:  "GET",
		Handler: fc.makeHandler(fc.QueryFarmPoolRewards, nil),
	})

	endpoints = append(endpoints, kit.Endpoint{
		URI:     "/farm/infos",
		Method:  "GET",
		Handler: fc.makeHandler(fc.QueryTotalVolumeLock, nil),
	})

	endpoints = append(endpoints, kit.Endpoint{
		URI:     "/farm/claim",
		Method:  "POST",
		Handler: fc.makeHandler(fc.Claim, &types.ClaimRequest{}),
	})

	endpoints = append(endpoints, kit.Endpoint{
		URI:     "/farm/claim",
		Method:  "GET",
		Handler: fc.makeHandler(fc.QueryClaim, nil),
	})

	endpoints = append(endpoints, kit.Endpoint{
		URI:     "/farm/coins_price",
		Method:  "GET",
		Handler: fc.makeHandler(fc.QueryCoinPrice, nil),
	})

	endpoints = append(endpoints, kit.Endpoint{
		URI:     "/farm/cache",
		Method:  "DELETE",
		Handler: fc.makeHandler(fc.DeleteCache, nil),
	})
	return endpoints
}
func (fc farmController) DeleteCache(ctx context.Context, _ interface{}) (interface{}, error) {

	redisKey, _ := fc.GetStringValue(ctx, "key")
	err := fc.DelCache(redisKey)
	if err != nil {
		return nil, err
	}
	return struct{}{}, nil
}

func (fc farmController) QueryCoinPrice(ctx context.Context, _ interface{}) (interface{}, error) {

	coins, err := fc.QueryAllPrice()
	if err != nil {
		logger.Error("Query coin list failed", logger.String("error", err.Error()))
		return nil, err
	}
	return types.QueryCoinPriceListResponse{Coins: coins}, nil
}

func (fc farmController) QueryCoinList(ctx context.Context, _ interface{}) (interface{}, error) {
	coins, err := fc.QueryAllCoins()
	if err != nil {
		logger.Error("Query coin list failed", logger.String("error", err.Error()))
		return nil, err
	}
	return types.QueryCoinListResponse{Coins: coins}, nil
}

func (fc farmController) QueryFarmPoolRewards(ctx context.Context, _ interface{}) (interface{}, error) {
	address, _ := fc.GetStringValue(ctx, "address")
	if len(address) == 0 {
		return types.QueryFarmPoolRewardsResponse{}, nil
	}
	if err := sdktype.ValidateAccAddress(address); err != nil {
		return nil, errors.Errorf("invalid address: %s ", address)
	}
	poolId, _ := fc.GetStringValue(ctx, "pool_id")
	pools, height, err := fc.PoolService.QueryAllRewards(address, poolId)
	if err != nil {
		logger.Error("Query farm pool list failed", logger.Any("error", err))
		return nil, err
	}
	return types.QueryFarmPoolRewardsResponse{
		Pools:  pools,
		Height: height,
	}, nil
}

func (fc farmController) QueryFarmPoolList(ctx context.Context, _ interface{}) (interface{}, error) {
	pools, err := fc.PoolService.QueryAllPools()
	if err != nil {
		logger.Error("Query farm pool list failed", logger.Any("error", err))
		return nil, err
	}
	return types.QueryFarmPoolListResponse{
		Pools: pools,
	}, nil
}

func (fc farmController) QueryTotalVolumeLock(_ context.Context, _ interface{}) (interface{}, error) {
	volumeStr, err := fc.PoolService.QueryTotalVolumeLock()
	if err != nil {
		logger.Error("Query total volume failed", logger.String("error", err.Error()))
		return nil, err
	}

	return types.QueryTotalVolumeLockResponse{TotalVolumeLocked: volumeStr}, nil
}

func (fc farmController) Claim(_ context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*types.ClaimRequest)
	if !ok {
		logger.Error("request args type error")
		return nil, errors.New("request args error")
	}
	if err := fc.PoolService.SaveClaim(req.Address); err != nil {
		return nil, err
	}

	return types.ClaimResponse{}, nil
}

func (fc farmController) QueryClaim(ctx context.Context, request interface{}) (interface{}, error) {
	address, _ := fc.GetStringValue(ctx, "address")
	if len(address) > 0 {
		if err := sdktype.ValidateAccAddress(address); err != nil {
			return nil, errors.Errorf("invalid address: %s ", address)
		}
	}

	has, err := fc.PoolService.CheckClaim(address)
	if err != nil {
		return nil, err
	}
	return types.QueryClaimResponse{Claimed: has}, nil
}

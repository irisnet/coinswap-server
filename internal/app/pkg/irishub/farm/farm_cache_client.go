package farm

import (
	"context"
	"fmt"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"time"
)

const (
	RedisKeyFarmer           = "farm:%s"
	RedisKeyFarmerExpiration = 5 * time.Second
)

var _ Client = &LcdFarmCacheClient{}

type LcdFarmCacheClient LcdFarmClient

func NewLcdFarmCacheClient(url string) Client {
	return &LcdFarmCacheClient{
		BaseUrl: url,
	}
}

func (lc *LcdFarmCacheClient) FarmPools(ctx context.Context, req *QueryFarmPoolsRequest) (*QueryFarmPoolsResponse, error) {

	resp, err := (*LcdFarmClient)(lc).FarmPools(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (lc *LcdFarmCacheClient) FarmPool(ctx context.Context, req *QueryFarmPoolRequest) (*QueryFarmPoolResponse, error) {

	resp, err := (*LcdFarmClient)(lc).FarmPool(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (lc *LcdFarmCacheClient) Farmer(ctx context.Context, req QueryFarmerRequest) (*QueryFarmerResponse, error) {

	var resp QueryFarmerResponse
	err := redis.GetObject(fmt.Sprintf(RedisKeyFarmer, req.Farmer), &resp)
	if err == nil {
		return &resp, nil
	}

	farmerResp, err := (*LcdFarmClient)(lc).Farmer(ctx, req)

	if err != nil {
		return nil, err
	}

	if err := redis.SetObject(fmt.Sprintf(RedisKeyFarmer, req.Farmer), farmerResp, RedisKeyFarmerExpiration); err != nil {
		logger.Error("Save  farmer to redis failed", logger.String("error", err.Error()))
	}

	return farmerResp, nil
}

func (lc *LcdFarmCacheClient) Params(ctx context.Context, req QueryParamsRequest) (*QueryParamsResponse, error) {

	resp, err := (*LcdFarmClient)(lc).Params(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

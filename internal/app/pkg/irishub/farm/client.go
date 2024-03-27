package farm

import "context"

type Client interface {
	FarmPools(ctx context.Context, req *QueryFarmPoolsRequest) (*QueryFarmPoolsResponse, error)

	FarmPool(ctx context.Context, req *QueryFarmPoolRequest) (*QueryFarmPoolResponse, error)

	Farmer(ctx context.Context, req QueryFarmerRequest) (*QueryFarmerResponse, error)

	Params(ctx context.Context, req QueryParamsRequest) (*QueryParamsResponse, error)
}

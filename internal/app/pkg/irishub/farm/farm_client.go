package farm

import (
	"context"
	"encoding/json"
	"fmt"
	util "github.com/irisnet/coinswap-server/internal/app/server/types"
	"github.com/irisnet/irishub-sdk-go/types"
)

const (
	farmPools = "/irismod/farm/pools?pagination.offset=%v&pagination.limit=%v&pagination.count_total=%v"
	farmPool  = "/irismod/farm/pool/%s"    //irismod/farm/pool/{id}
	farmer    = "/irismod/farm/farmers/%s" //"irismod/farm/farmers/{farmer}/
	params    = "/irismod/farm/params"
)

type (
	LcdFarmClient struct {
		BaseUrl string
	}

	QueryFarmPoolsRequest struct {
		Pagination *types.PageRequest `json:"pagination,omitempty"`
	}

	QueryFarmPoolRequest struct {
		Name string `json:"name,omitempty"`
	}

	QueryFarmPoolResponse struct {
		Pool *FarmPoolEntry `json:"pool,omitempty"`
	}

	PageResponse struct {
		NextKey []byte `json:"next_key,omitempty"`
		Total   string `json:"total,omitempty"`
	}
	QueryFarmPoolsResponse struct {
		Pools      []*FarmPoolEntry `json:"pools,omitempty"`
		Pagination *PageResponse    `json:"pagination,omitempty"`
	}

	QueryFarmerRequest struct {
		Farmer   string `json:"farmer,omitempty"`
		PoolName string `json:"pool_name,omitempty"`
	}

	QueryFarmerResponse struct {
		List   []*LockedInfo `json:"list"`
		Height string        `json:"height"`
	}

	LockedInfo struct {
		PoolId        string      `json:"pool_id"`
		Locked        types.Coin  `json:"locked"`
		PendingReward types.Coins `json:"pending_reward"`
	}

	QueryParamsRequest struct {
	}

	QueryParamsResponse struct {
		Params Params `json:"params"`
	}

	Params struct {
		CreatePoolFee       types.Coin `json:"create_pool_fee"`
		MaxRewardCategories uint32     `json:"max_reward_categories,omitempty"`
	}

	FarmPoolEntry struct {
		Id              string       `json:"id,omitempty"`
		Creator         string       `json:"creator,omitempty"`
		Description     string       `json:"description,omitempty"`
		StartHeight     string       `json:"start_height,omitempty"`
		EndHeight       string       `json:"end_height,omitempty"`
		Editable        bool         `json:"editable,omitempty"`
		Expired         bool         `json:"expired,omitempty"`
		TotalLptLocked  types.Coin   `json:"total_lpt_locked"`
		TotalReward     []types.Coin `json:"total_reward"`
		RemainingReward []types.Coin `json:"remaining_reward"`
		RewardPerBlock  []types.Coin `json:"reward_per_block"`
	}
)

func NewLcdFarmClient(baseUrl string) Client {
	return &LcdFarmClient{BaseUrl: baseUrl}
}

func (farm *LcdFarmClient) FarmPools(ctx context.Context, req *QueryFarmPoolsRequest) (*QueryFarmPoolsResponse, error) {

	var resp QueryFarmPoolsResponse

	if err := farm.get(fmt.Sprintf(farmPools, req.Pagination.Offset, req.Pagination.Limit, req.Pagination.CountTotal), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (farm *LcdFarmClient) FarmPool(ctx context.Context, req *QueryFarmPoolRequest) (*QueryFarmPoolResponse, error) {

	var resp QueryFarmPoolResponse

	if err := farm.get(fmt.Sprintf(farmPool, req.Name), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (farm *LcdFarmClient) Farmer(ctx context.Context, req QueryFarmerRequest) (*QueryFarmerResponse, error) {

	var resp QueryFarmerResponse

	if err := farm.get(fmt.Sprintf(farmer, req.Farmer), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (farm *LcdFarmClient) Params(ctx context.Context, req QueryParamsRequest) (*QueryParamsResponse, error) {

	var resp QueryParamsResponse

	if err := farm.get(params, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (farm *LcdFarmClient) get(url string, dest interface{}) error {

	respBytes, err := util.Get(farm.BaseUrl + url)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(respBytes, dest); err != nil {
		return err
	}
	return nil
}

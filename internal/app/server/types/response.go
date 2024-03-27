package types

import (
	"encoding/json"
	sdktype "github.com/irisnet/irishub-sdk-go/types"
)

type (
	CoinItem struct {
		Name             string        `json:"name"`
		Symbol           string        `json:"symbol"`
		Denom            string        `json:"denom"`
		Scale            int           `json:"scale"`
		CoinID           string        `json:"coin_id"`
		Offline          bool          `json:"offline"`
		Platform         string        `json:"platform"`
		Protocol         int           `json:"protocol"`
		LpToken          string        `json:"lp_token"`
		Icon             string        `json:"icon"`
		Tips             string        `json:"tips"`
		IbcTransferInfos []IbcInfoItem `json:"ibc_transfer_infos"`
	}
	CoinPrice struct {
		Denom string `json:"denom"`
		Price string `json:"price"`
	}
	IbcInfoItem struct {
		WithdrawInfo InPath      `json:"withdraw_info"`
		DepositInfo  []TraceInfo `json:"deposit_info"`
	}
	TraceInfo struct {
		Port      string `json:"port"`
		Channel   string `json:"channel"`
		ChainId   string `json:"chain_id"`
		ChainName string `json:"chain_name"`
		Denom     string `json:"denom"`
		Platform  string `json:"platform"`
	}
	InPath struct {
		Port      string `json:"port"`
		Channel   string `json:"channel"`
		ChainId   string `json:"chain_id"`
		ChainName string `json:"chain_name"`
	}

	Prices []CoinPrice
	Coins  []CoinItem

	Pool struct {
		ID           string  `json:"id"`
		Code         string  `json:"code"`
		Icon         string  `json:"icon"`
		Name         string  `json:"name"`
		VolumeLocked float64 `json:"volume_locked"`
		Reward       string  `json:"reward"`
		APY          float64 `json:"apy"`
		Status       int     `json:"status"`
		ActivityOn   bool    `json:"activity_on"`
	}
	PoolReward struct {
		PoolId string    `json:"pool_id"`
		Staked uint64    `json:"staked"`
		Earned []CoinStr `json:"earned"`
	}
	CoinStr struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}

	ClaimRequest struct {
		Address string `validate:"required" json:"address"`
	}

	ClaimResponse struct{}

	QueryClaimResponse struct {
		Claimed bool `json:"claimed"`
	}

	QueryCoinListResponse struct {
		Coins []CoinItem `json:"coins"`
	}
	QueryCoinPriceListResponse struct {
		Coins []CoinPrice `json:"coins"`
	}

	QueryTotalVolumeLockResponse struct {
		TotalVolumeLocked string `json:"volume_total_locked"`
	}

	QueryFarmPoolListResponse struct {
		Pools []Pool `json:"pools"`
	}
	QueryFarmPoolRewardsResponse struct {
		Height string       `json:"height"`
		Pools  []PoolReward `json:"pools"`
	}

	QueryLiquidityPoolResponse struct {
		Denom     string       `json:"denom"`
		Standard  sdktype.Coin `json:"standard"`
		Token     sdktype.Coin `json:"token"`
		Liquidity sdktype.Coin `json:"liquidity"`
		Fee       string       `json:"fee"`
		UpdateAt  int64        `json:"update_at"`
	}
	QueryNotice struct {
		Notice string `json:"notice"`
	}
)

func (cs Coins) MarshalBinary() ([]byte, error) {
	return json.Marshal(cs)
}

func (p Pool) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

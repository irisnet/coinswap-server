package model

import "github.com/irisnet/coinswap-server/internal/app/server/types"

type IWhiteList interface {
	FindAll() ([]WhiteList, error)
}

type ILpTokens interface {
	FindAll() ([]LpTokens, error)
}

type IAsset interface {
	QueryAllScale() ([]Asset, error)
	FindAll() ([]Asset, error)
	FindCoinByDenom(denom string) (Asset, error)
}

type IStatisticInfo interface {
	FindLiquidityPoolByDenom(denom string) (types.QueryLiquidityPoolResponse, error)
	FindLiquidityPoolByDenoms(denoms []string) ([]types.QueryLiquidityPoolResponse, error)
	FindTotalVolumeLock() (string, error)
}

type IClaimAddress interface {
	Exist(address string) (bool, error)
	Save(info ClaimAddress) error
}

package service

import (
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	"strings"
)

type LiquidityService struct {
}

func NewLiquidityService() LiquidityService {
	return LiquidityService{}
}
func (ls *LiquidityService) QueryPool(denom string) ([]types.QueryLiquidityPoolResponse, error) {
	denoms := strings.Split(denom, ",")
	response, err := statisticModel.FindLiquidityPoolByDenoms(denoms)
	if err != nil {
		return response, err
	}

	return response, nil
}

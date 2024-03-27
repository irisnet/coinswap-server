package controller

import (
	"context"
	"github.com/irisnet/coinswap-server/internal/app/pkg/kit"
	"github.com/irisnet/coinswap-server/internal/app/server/service"
)

type liquidityController struct {
	BaseController
	ls service.LiquidityService
}

func NewLiquidityController(bc BaseController, ls service.LiquidityService) liquidityController {
	return liquidityController{bc, ls}
}

func (lc liquidityController) GetEndpoints() (endpoints []kit.Endpoint) {
	endpoints = append(endpoints, kit.Endpoint{
		URI:     "/liquidity/pool",
		Method:  "GET",
		Handler: lc.makeHandler(lc.QueryPool, nil),
	})
	return endpoints
}

func (lc liquidityController) QueryPool(ctx context.Context, _ interface{}) (interface{}, error) {
	denom, err := lc.GetStringValue(ctx, "denom")
	if err != nil {
		return nil, err
	}
	return lc.ls.QueryPool(denom)
}

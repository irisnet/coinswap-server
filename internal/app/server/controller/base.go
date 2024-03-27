package controller

import (
	"context"
	"github.com/irisnet/coinswap-server/internal/app/pkg/kit"
	"github.com/irisnet/coinswap-server/internal/app/server/monitor"
	"github.com/irisnet/coinswap-server/internal/app/server/service"
	"github.com/irisnet/coinswap-server/internal/app/server/task"
)

type (
	//BaseController define a base controller for all http Controller
	BaseController struct {
		kit.Controller
	}
)

// return all the controllers of the app server
func GetAllControllers() []kit.IController {
	bc := BaseController{
		Controller: kit.NewController(),
	}
	coinService := service.NewCoinService()
	metricsService := monitor.NewMetricsService()
	poolService := service.NewPoolService(coinService)
	liquidityService := service.NewLiquidityService()
	noticeService := service.NewNoticeService()

	controllers := []kit.IController{
		farmController{bc, coinService, poolService},
		NewMetricsController(metricsService),
		NewLiquidityController(bc, liquidityService),
		NewNoticeController(bc, noticeService),
	}

	task.Start()
	return controllers
}

// makeHandler create a http hander for request
func (bc BaseController) makeHandler(h kit.Handler, request interface{}) *kit.Server {
	return bc.MakeHandler(
		bc.wrapHandler(h),
		request,
		[]kit.RequestFunc{},
		nil,
		[]kit.ServerResponseFunc{},
	)
}

func (bc BaseController) wrapHandler(h kit.Handler) kit.Handler {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		resp, err := h(ctx, request)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

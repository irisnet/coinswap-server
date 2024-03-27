package server

import (
	"github.com/irisnet/coinswap-server/config"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/pkg/kit"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/controller"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
)

type FarmServerApp struct {
}

func NewFarmServerApp() kit.Application {
	return FarmServerApp{}
}

// GetEndpoints return all the endpoints for http server
func (fsa FarmServerApp) GetEndpoints() []kit.Endpoint {
	var rs []kit.Endpoint

	ctls := controller.GetAllControllers()
	for _, c := range ctls {
		rs = append(rs, c.GetEndpoints()...)
	}
	return rs
}

func (fsa FarmServerApp) Initialize() {
	conf := config.Get()
	model.Init(&conf)
	//初始化redis连接
	redis.Connect(conf.Redis.Address, conf.Redis.Password, conf.Redis.DB)
	//初始化irishub连接
	irishub.Connect(
		conf.IRIShub.RcpAddr,
		conf.IRIShub.GrpcAddr,
		conf.IRIShub.ChainID,
		conf.IRIShub.Fee,
	)
}

func (fsa FarmServerApp) Stop() {
	model.Close()
}

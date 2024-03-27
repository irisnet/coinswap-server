package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/pkg/redis"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	util "github.com/irisnet/coinswap-server/internal/app/server/types"
	"sync"
	"time"
)

var (
	mutex           sync.Mutex
	_cronTaskStatus map[string]int
)

type MetricsService struct {
}

func NewMetricsService() MetricsService {
	return MetricsService{}
}

func (ms *MetricsService) QueryIrishubConnectionStatus() int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := irishub.GetClient().Status(ctx)
	if err != nil {
		logger.Error("rpc node connection exception, err is " + err.Error())
		return NodeStatusNotReachable

	}
	return NodeStatusReachable
}

func (ms *MetricsService) QueryRedisConnectionStatus() int {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := redis.Info(ctx, fmt.Sprint(model.GetConf().Redis.DB))
	if err != nil {
		logger.Error("redis connection exception, err is " + err.Error())
		return NodeStatusNotReachable
	}
	return NodeStatusReachable
}

func (ms *MetricsService) QueryLcdConnectionStatus() int {
	respBytes, err := util.Get(model.GetConf().IRIShub.LcdAddr + "/cosmos/base/tendermint/v1beta1/node_info")
	if err != nil {
		logger.Error("lcd connection exception, err is " + err.Error())
		return 0
	}
	var data LcdNodeInfoResp
	if err := json.Unmarshal(respBytes, &data); err != nil {
		logger.Error("Lcd NodeInfo Resp Unmarshal err is " + err.Error())
	}
	if data.NodeInfo.Network != "" {
		return 1
	}
	return 0
}

func (ms *MetricsService) QueryCronTaskStatus() map[string]int {
	return _cronTaskStatus
}

func SetCronTaskStatus(taskName string, status int) {
	mutex.Lock()
	if len(_cronTaskStatus) == 0 {
		_cronTaskStatus = make(map[string]int, 1)
	}
	_cronTaskStatus[taskName] = status
	mutex.Unlock()
}

package controller

import (
	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/irisnet/coinswap-server/internal/app/pkg/kit"
	"github.com/irisnet/coinswap-server/internal/app/server/monitor"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

const namespace = "dapp_server"

type metricsController struct {
	nodeConnectionStatus metrics.Gauge
	lcdConnectionStatus  metrics.Gauge
	redisStatus          metrics.Gauge
	cronTaskStatus       metrics.Gauge
	next                 http.Handler
	ms                   monitor.MetricsService
}

func NewMetricsController(ms monitor.MetricsService) kit.IController {
	irishubNodeConnectionStatus := kitprometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "irishub",
		Name:      "node_connection_status",
		Help:      "irishub_node_connection_status node connection status of irishub from dapp-server",
	}, []string{})
	irishubNodeLcdConnectionStatus := kitprometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "irishub",
		Name:      "node_lcd_connection_status",
		Help:      "irishub_node_lcd_connection_status node lcd connection status (1:Reachable  0:NotReachable)",
	}, []string{})
	redisNodeConnectionStatus := kitprometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "redis",
		Name:      "node_connection_status",
		Help:      "redis_node_connection_status node connection status of redis service from dapp-server",
	}, []string{})

	cronTaskStatus := kitprometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "cron_task",
		Name:      "working_status",
		Help:      "dapp_server_cron_task_working_status dapp-server cron task working status",
	}, []string{"taskname"})

	return metricsController{
		nodeConnectionStatus: irishubNodeConnectionStatus,
		lcdConnectionStatus:  irishubNodeLcdConnectionStatus,
		redisStatus:          redisNodeConnectionStatus,
		cronTaskStatus:       cronTaskStatus,
		next:                 promhttp.Handler(),
		ms:                   ms,
	}
}

// GetRouters implement the method GetRouter of the interface IController
func (mc metricsController) GetEndpoints() []kit.Endpoint {
	var ends []kit.Endpoint
	ends = append(ends, kit.Endpoint{
		URI:     "/metrics",
		Method:  "GET",
		Handler: mc,
	})
	return ends
}

func (mc metricsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wg := &sync.WaitGroup{}
	wg.Add(4)
	mc.IrishubNodeStatus(wg)
	mc.IrisNodeLcdStatus(wg)
	mc.RedisNodeStatus(wg)
	mc.CronTaskStatus(wg)
	wg.Wait()
	mc.next.ServeHTTP(w, r)
}

func (mc metricsController) IrishubNodeStatus(l *sync.WaitGroup) {
	defer l.Done()
	value := mc.ms.QueryIrishubConnectionStatus()
	mc.nodeConnectionStatus.Set(float64(value))
}

func (mc metricsController) IrisNodeLcdStatus(l *sync.WaitGroup) {
	defer l.Done()
	value := mc.ms.QueryLcdConnectionStatus()
	mc.lcdConnectionStatus.Set(float64(value))
}

func (mc metricsController) RedisNodeStatus(l *sync.WaitGroup) {
	defer l.Done()
	value := mc.ms.QueryRedisConnectionStatus()
	mc.redisStatus.Set(float64(value))
}

func (mc metricsController) CronTaskStatus(l *sync.WaitGroup) {
	defer l.Done()
	values := mc.ms.QueryCronTaskStatus()
	for key, val := range values {
		mc.cronTaskStatus.With("taskname", key).Set(float64(val))
	}
}

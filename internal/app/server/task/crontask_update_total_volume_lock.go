package task

import (
	"fmt"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"github.com/irisnet/coinswap-server/internal/app/server/monitor"
	"github.com/irisnet/coinswap-server/internal/app/server/service"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	sdktype "github.com/irisnet/irishub-sdk-go/types"
	"github.com/pkg/errors"
	"github.com/qiniu/qmgo"
	"time"
)

func init() {
	RegisterTasks(&SyncUpdateTotalVolumnLockTask{})
}

type SyncUpdateTotalVolumnLockTask struct {
	service.CoinService
}

func (s SyncUpdateTotalVolumnLockTask) Name() string {
	return types.SyncUpdateTotalVolumnLockTaskName
}

func (s SyncUpdateTotalVolumnLockTask) Cron() int {
	if model.GetConf().Task.CronTimeUpdateLiquidityPool > 0 {
		return model.GetConf().Task.CronTimeUpdateTotalVolumeLock
	}
	return types.CronTimeUpdateTotalVolumnLock
}

func (s SyncUpdateTotalVolumnLockTask) Start() {
	timeInterval := s.Cron()
	RunTimer(timeInterval, Sec, func() {
		monitor.SetCronTaskStatus(types.SyncUpdateTotalVolumnLockTaskName, 0)
		if err := s.DoTask(nil); err != nil {
			logger.Error(err.Error(), logger.String("taskName", s.Name()))
		}
	})

}

func SaveOrUpdateStatisticInfo(value string, statisticName string) error {
	if statisticName == "" || value == "" {
		return fmt.Errorf("invalid value or statisticName")
	}
	var statisticModel model.StatisticInfo
	data, err := statisticModel.FindByName(statisticName)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			if err := statisticModel.Save(model.StatisticInfo{
				Name:          statisticName,
				StatisticData: value,
				CreateAt:      time.Now().Unix(),
				UpdateAt:      time.Now().Unix(),
			}); err != nil {
				//logger.Error("Save total volume lock", logger.String("error", err.Error()))
				return err
			}
			return nil
		}

		//logger.Error("Find total volume lock", logger.String("error", err.Error()))
		return err
	}
	data.StatisticData = value
	if err := statisticModel.Update(data); err != nil {
		//logger.Error("Update total volume lock", logger.String("error", err.Error()))
		return err
	}
	return nil
}

func (s SyncUpdateTotalVolumnLockTask) DoTask(fn func(string) chan bool) error {
	if !model.GetConf().Task.Enable {
		logger.Warn("CronTask Config Disable", logger.String("taskName", s.Name()))
		return nil
	}
	volume, err := s.calculateTotalVolumeLock()
	if err != nil {
		monitor.SetCronTaskStatus(types.SyncUpdateTotalVolumnLockTaskName, -1)
		return errors.Wrap(err, "calculate TotalVolumeLock failed")
	}
	if err := SaveOrUpdateStatisticInfo(volume, types.StatisticTVL); err != nil {
		monitor.SetCronTaskStatus(types.SyncUpdateTotalVolumnLockTaskName, -1)
		return errors.Wrap(err, "Update or save total volume lock failed")
	}
	monitor.SetCronTaskStatus(types.SyncUpdateTotalVolumnLockTaskName, 1)

	return nil
}

func (s SyncUpdateTotalVolumnLockTask) calculateTotalVolumeLock() (string, error) {
	result, err := irishub.Coinswap().QueryAllPools(sdktype.PageRequest{
		//this no use page and count
	})
	if err != nil {
		//logger.Error("Query liquidity pool failed", logger.String("error", err.Error()))
		return "", err
	}

	farmLpAmtMap, err := s.caculateFarmLPAmt()
	if err != nil {
		logger.Error("caculate Farm LP Amt failed," + err.Error())
	}

	totalBaseAmt := sdktype.ZeroDec()
	for _, pool := range result.Pools {

		if len(farmLpAmtMap) > 0 {
			if lpAmt, ok := farmLpAmtMap[pool.Lpt.Denom]; ok && !pool.Lpt.Amount.IsZero() {
				lpIrisAmt := lpAmt.Quo(sdktype.NewDecFromInt(pool.Lpt.Amount)).Mul(sdktype.NewDecFromInt(pool.Standard.Amount))
				totalBaseAmt = totalBaseAmt.Add(lpIrisAmt)
			}
		}

		totalBaseAmt = totalBaseAmt.Add(sdktype.NewDecFromInt(pool.Standard.Amount))
		//logger.Debug("Query liquidity Pool status", logger.String("baseToken", pool.Standard.String()),
		//	logger.String("Token", pool.Token.String()),
		//	logger.String("liquidity", pool.Lpt.String()),
		//)
	}
	_, volumeByUSD, err := s.PriceIRISByUSD(totalBaseAmt.MulInt64(2))
	if err != nil {
		return "", err
	}
	//logger.Debug("Calculate total volume locked  by usd", logger.String("volume", volumeByUSD.String()))
	return volumeByUSD.String(), nil
}

func (s SyncUpdateTotalVolumnLockTask) caculateFarmLPAmt() (map[string]sdktype.Dec, error) {
	var (
		offset int64 = 0
		limit  int64 = 10
	)
	totalLpAmt := make(map[string]sdktype.Dec, 1)
	for {
		farms, err := new(model.Farms).FindLPTokenByPageSize(offset, limit)
		if err != nil {
			return nil, err
		}
		for _, val := range farms {
			amt, err := sdktype.NewDecFromStr(val.TotalLptLocked.Amount)
			if err != nil {
				logger.Warn(err.Error(),
					logger.String("LpTokenDenom", val.TotalLptLocked.Denom),
					logger.String("LpTokenAmt", val.TotalLptLocked.Amount),
					logger.String("pool_id", val.PoolId))
				continue
			}
			data, ok := totalLpAmt[val.TotalLptLocked.Denom]
			if ok {
				totalLpAmt[val.TotalLptLocked.Denom] = data.Add(amt)
			} else {
				totalLpAmt[val.TotalLptLocked.Denom] = amt
			}
		}

		if len(farms) < int(limit) {
			break
		}
		offset += limit
	}

	return totalLpAmt, nil
}

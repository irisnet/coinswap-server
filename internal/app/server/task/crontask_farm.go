package task

import (
	"context"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub"
	"github.com/irisnet/coinswap-server/internal/app/pkg/irishub/farm"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/server/model"
	"github.com/irisnet/coinswap-server/internal/app/server/model/cache"
	"github.com/irisnet/coinswap-server/internal/app/server/monitor"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	sdktype "github.com/irisnet/irishub-sdk-go/types"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

func init() {
	RegisterTasks(&SyncFarmTask{})
}

type SyncFarmTask struct {
}

func (t *SyncFarmTask) Name() string {
	return types.SyncFarmTaskName
}

func (t *SyncFarmTask) Cron() int {
	return types.CronTimeSyncFarmTask
}

func (t *SyncFarmTask) DoTask(fn func(string) chan bool) error {
	return nil
}

func (t *SyncFarmTask) Start() {

	RunTimer(t.Cron(), Sec, func() {
		monitor.SetCronTaskStatus(types.SyncFarmTaskName, 0)
		if !model.GetConf().Task.Enable {
			logger.Warn("CronTask Config Disable", logger.String("taskName", t.Name()))
			return
		}
		//logger.Info("#### syncFarmTask start ###")
		err := NewFarmPoolSyncer().Execute()
		if err != nil {
			monitor.SetCronTaskStatus(types.SyncFarmTaskName, -1)
			logger.Error("execute task occurs err", logger.String("taskName", t.Name()), logger.String("err", err.Error()))
			return
		}
		monitor.SetCronTaskStatus(types.SyncFarmTaskName, 1)
		//logger.Info("### syncFarmTask  finished ###")
	})
}

type FarmPoolSyncer struct {
	whiteMap      map[string]model.WhiteList
	farmInDb      map[string]model.Farms
	farmInChain   map[string]model.Farms
	farmForAdd    []model.Farms
	farmForUpdate []model.Farms
}

func NewFarmPoolSyncer() *FarmPoolSyncer {
	return &FarmPoolSyncer{
		whiteMap:    make(map[string]model.WhiteList),
		farmInDb:    make(map[string]model.Farms),
		farmInChain: make(map[string]model.Farms),
	}
}

func (fps *FarmPoolSyncer) Execute() error {
	if err := fps.LoadDataFromDb(); err != nil {
		monitor.SetCronTaskStatus(types.SyncFarmTaskName, -1)
		return err
	}

	if err := fps.MergeFarmFromChain(); err != nil {
		monitor.SetCronTaskStatus(types.SyncFarmTaskName, -1)
		return err
	}

	//return fps.UpdateVisibleForLocalFarm()
	return nil

}

// LoadDataFromDb load white list and farm from db ,for comparing with the farm from chain later
func (fps *FarmPoolSyncer) LoadDataFromDb() error {

	if err := fps.FindWhiteListFromDb(); err != nil {
		return err
	}

	return fps.FindFarmPoolFromDb()
}

// FindWhiteListFromDb get white list from db
func (fps *FarmPoolSyncer) FindWhiteListFromDb() error {
	whiteList, err := new(cache.WhiteListCache).FindAll()
	if err != nil {
		return errors.Wrapf(err, "whiteList.FindAll()")
	}

	for _, wl := range whiteList {
		fps.whiteMap[wl.PoolId] = model.WhiteList{
			OrderBy:  wl.OrderBy,
			PoolName: wl.PoolName,
		}
	}
	return nil
}

// FindFarmPoolFromDb get farm from db
func (fps *FarmPoolSyncer) FindFarmPoolFromDb() error {

	farms, err := new(model.Farms).FindAll()
	if err != nil {
		return errors.Wrapf(err, "farms.FindAll()")
	}

	for _, farm := range farms {
		fps.farmInDb[farm.PoolId] = farm
	}
	return nil
}

// MergeFarmFromChain   save or update farm to db ,according to chain
func (fps *FarmPoolSyncer) MergeFarmFromChain() error {

	if err := fps.FindFarmPoolFromChain(); err != nil {
		return err
	}

	if err := fps.UpdateFarm(); err != nil {
		return err
	}

	return fps.AddFarm()
}

// FindFarmPoolFromChain get farm  from chain
func (fps *FarmPoolSyncer) FindFarmPoolFromChain() error {

	var (
		offset uint64 = 0
		limit  uint64 = 10
	)

	var farmPoolRequest = &farm.QueryFarmPoolsRequest{Pagination: &sdktype.PageRequest{
		Offset:     offset,
		Limit:      limit,
		CountTotal: true,
	}}

	// loop until reach the end
	for {
		farmPoolResponse, err := irishub.Farm().FarmPools(context.Background(), farmPoolRequest)
		if err != nil {
			return errors.Wrapf(err, "irishub.Farm.FarmPools(%v)", farmPoolRequest)
		}

		fps.PickFarmForUpdateOrAdd(farmPoolResponse.Pools)

		total, err := strconv.ParseUint(farmPoolResponse.Pagination.Total, 10, 64)
		if err != nil {
			return errors.Wrapf(err, "strconv.ParseUint(%v,10,64)", farmPoolResponse.Pagination.Total)
		}
		if total <= offset+limit {
			break
		}

		offset += limit
		farmPoolRequest.Pagination.Offset = offset
	}
	return nil
}

// PickFarmForUpdateOrAdd    divide the farms into 2 categories., One is to add and the other is to modify
func (fps *FarmPoolSyncer) PickFarmForUpdateOrAdd(farmPoolEntries []*farm.FarmPoolEntry) {

	for _, entry := range farmPoolEntries {
		var farmFromChain = NewFarmFromFarmPoolEntry(entry)

		// if it is in white list,fill the ctl info
		if ctlInfo, exist := fps.whiteMap[farmFromChain.PoolId]; exist {
			farmFromChain.OrderBy = ctlInfo.OrderBy
			farmFromChain.Name = ctlInfo.PoolName
			farmFromChain.Visible = true
		}

		// farm already exist in db
		dbFarm, existInDb := fps.farmInDb[farmFromChain.PoolId]
		if existInDb {
			// local farm ,should not update
			//if !dbFarm.IsLocalFarm {
			farmFromChain.Id = dbFarm.Id
			fps.farmForUpdate = append(fps.farmForUpdate, farmFromChain)

			//}
		} else {
			fps.farmForAdd = append(fps.farmForAdd, farmFromChain)
		}

	}
}

// UpdateFarm    update the farm in db
func (fps *FarmPoolSyncer) UpdateFarm() error {

	for _, farm := range fps.farmForUpdate {
		if err := farm.ReplaceOneByPoolID(); err != nil {
			return errors.Wrapf(err, "farm.ReplaceOneByPoolId()")
		}
	}

	return nil
}

// AddFarm add farm to db
func (fps *FarmPoolSyncer) AddFarm() error {

	if len(fps.farmForAdd) == 0 {
		return nil
	}

	err := new(model.Farms).InsertMany(fps.farmForAdd)
	if err != nil {
		return errors.Wrapf(err, "Farm.InsertMany(%v)", fps.farmForAdd)
	}
	return nil
}

////UpdateVisibleForLocalFarm  set visible  for the farms  which is local
//func (fps *FarmPoolSyncer) UpdateVisibleForLocalFarm() error {
//
//	// set visible for  the farm in db and not in chain
//	for poolId, dbFarm := range fps.farmInDb {
//
//		// only localFarm should  care white list
//		if !dbFarm.IsLocalFarm {
//			continue
//		}
//
//		if ctlInfo, exist := fps.whiteMap[poolId]; exist {
//			dbFarm.OrderBy = ctlInfo.OrderBy
//			dbFarm.Visible = true
//		} else {
//			dbFarm.OrderBy = ctlInfo.OrderBy
//			dbFarm.Visible = false
//		}
//
//		if err := dbFarm.ReplaceOneByPoolID(); err != nil {
//			return errors.Wrapf(err, "farm.ReplaceOneByPoolID(%v)", dbFarm.PoolId)
//		}
//
//	}
//	return nil
//}

// NewFarmFromFarmPoolEntry  init farm  from a farm pool entry
func NewFarmFromFarmPoolEntry(farm *farm.FarmPoolEntry) model.Farms {
	convertFarmCoinToLocalCoin := func(fcoins []sdktype.Coin) (localCoins []model.Coin) {
		for _, fCoin := range fcoins {
			localCoins = append(localCoins, model.Coin{
				Denom:  fCoin.Denom,
				Amount: fCoin.Amount.String(),
			})
		}
		return
	}

	return model.Farms{
		Id:          primitive.NewObjectID(),
		PoolId:      farm.Id,
		Creator:     farm.Creator,
		Description: farm.Description,
		StartHeight: farm.StartHeight,
		EndHeight:   farm.EndHeight,
		Editable:    farm.Editable,
		Expired:     farm.Expired,
		TotalLptLocked: model.Coin{
			Denom:  farm.TotalLptLocked.Denom,
			Amount: farm.TotalLptLocked.Amount.String(),
		},
		TotalReward:     convertFarmCoinToLocalCoin(farm.TotalReward),
		RemainingReward: convertFarmCoinToLocalCoin(farm.RemainingReward),
		RewardPerBlock:  convertFarmCoinToLocalCoin(farm.RewardPerBlock),
		Visible:         false, // false default
		//IsLocalFarm:     false, // from chain should not be local
		CreateAt: time.Now().Unix(),
		UpdateAt: time.Now().Unix(),
	}
}

package model

import (
	"encoding/json"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type StatisticInfo struct {
	Name          string `bson:"name"`
	StatisticData string `bson:"statistic_data"`
	CreateAt      int64  `bson:"create_at"`
	UpdateAt      int64  `bson:"update_at"`
}

const CollectionNameStatisticInfo = "statistic_info"

func (d StatisticInfo) TableName() string {
	return CollectionNameStatisticInfo
}

func (d StatisticInfo) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:        []string{"-name"},
		Unique:     true,
		Background: true,
	})
	ensureIndexes(d.TableName(), indexes)
}

func (d StatisticInfo) PkKvPair() map[string]interface{} {
	return bson.M{"name": d.Name}
}

func (d StatisticInfo) FindByName(name string) (StatisticInfo, error) {
	var info StatisticInfo
	query := bson.M{
		"name": name,
	}
	fn := func(c *qmgo.Collection) error {
		return c.Find(_ctx, query).One(&info)
	}

	err := ExecCollection(d.TableName(), fn)

	return info, err
}
func (d StatisticInfo) FindLiquidityPoolByDenoms(denoms []string) ([]types.QueryLiquidityPoolResponse, error) {
	var res []types.QueryLiquidityPoolResponse
	info, err := d.FindByName(types.StatisticLiquidityPool)
	if err != nil {
		return []types.QueryLiquidityPoolResponse{}, err
	}
	var poolSlice []types.QueryLiquidityPoolResponse
	if err := json.Unmarshal([]byte(info.StatisticData), &poolSlice); err != nil {
		return []types.QueryLiquidityPoolResponse{}, err
	}

	denomMap := make(map[string]bool, len(denoms))
	for _, val := range denoms {
		denomMap[val] = true
	}

	for _, val := range poolSlice {
		if _, ok := denomMap[val.Denom]; ok {
			res = append(res, val)
		}
	}
	return res, err
}
func (d StatisticInfo) FindLiquidityPoolByDenom(denom string) (types.QueryLiquidityPoolResponse, error) {
	info, err := d.FindByName(types.StatisticLiquidityPool)
	if err != nil {
		return types.QueryLiquidityPoolResponse{}, err
	}
	var poolSlice []types.QueryLiquidityPoolResponse
	if err := json.Unmarshal([]byte(info.StatisticData), &poolSlice); err != nil {
		return types.QueryLiquidityPoolResponse{}, err
	}

	for _, val := range poolSlice {
		if val.Denom == denom {
			return val, nil
		}
	}
	return types.QueryLiquidityPoolResponse{}, err
}

func (d StatisticInfo) FindTotalVolumeLock() (string, error) {
	ret, err := d.FindByName(types.StatisticTVL)
	if err != nil {
		return "", err
	}
	return ret.StatisticData, err
}

func (d StatisticInfo) Save(info StatisticInfo) error {
	fn := func(c *qmgo.Collection) error {
		_, err := c.InsertOne(_ctx, info)
		if err != nil {
			return err
		}
		return nil
	}

	return ExecCollection(d.TableName(), fn)
}

func (d StatisticInfo) Update(info StatisticInfo) error {
	fn := func(c *qmgo.Collection) error {
		return c.UpdateOne(_ctx, info.PkKvPair(), bson.M{"$set": bson.M{
			"statistic_data": info.StatisticData,
			"update_at":      time.Now().Unix(),
		}})
	}

	return ExecCollection(d.TableName(), fn)
}

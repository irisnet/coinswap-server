package model

import (
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
)

type WhiteList struct {
	PoolId   string `bson:"pool_id"`
	PoolName string `bson:"pool_name"`
	OrderBy  int    `bson:"order_by"`
	CreateAt int64  `bson:"create_at" json:"-"`
	UpdateAt int64  `bson:"update_at" json:"-"`
}

const CollectionNameWhiteList = "white_list"

func (d WhiteList) TableName() string {
	return CollectionNameWhiteList
}

func (d WhiteList) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:        []string{"-pool_id"},
		Unique:     true,
		Background: true,
	})
	ensureIndexes(d.TableName(), indexes)
}

func (d WhiteList) PkKvPair() map[string]interface{} {
	return bson.M{"pool_id": d.PoolId}
}

func (w WhiteList) FindAll() (whiteList []WhiteList, err error) {

	err = NewCollection(w).Find(_ctx, bson.M{}).All(&whiteList)
	return
}

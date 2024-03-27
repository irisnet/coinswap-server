package model

import (
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
)

// LpTokens is an object representing the database table.
type LpTokens struct {
	Denom      string `bson:"denom"`
	Icon       string `bson:"icon"`
	CreateTime int64  `bson:"create_time"`
	UpdateTime int64  `bson:"update_time"`
}

const CollectionNameLpTokens = "lp_tokens"

func (d LpTokens) TableName() string {
	return CollectionNameLpTokens
}

func (d LpTokens) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:        []string{"-denom"},
		Unique:     true,
		Background: true,
	})
	ensureIndexes(d.TableName(), indexes)
}

func (d LpTokens) PkKvPair() map[string]interface{} {
	return bson.M{"denom": d.Denom}
}

func (d LpTokens) FindAll() ([]LpTokens, error) {
	var lpTokens []LpTokens
	query := bson.M{}
	fn := func(c *qmgo.Collection) error {
		return c.Find(_ctx, query).All(&lpTokens)
	}

	err := ExecCollection(d.TableName(), fn)
	return lpTokens, err
}

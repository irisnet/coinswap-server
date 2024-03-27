package model

import (
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
)

// ClaimAddress is an object representing the database table.
type ClaimAddress struct {
	Address    string `bson:"address"`
	CreateTime int64  `bson:"create_time"`
	UpdateTime int64  `bson:"update_time"`
}

const CollectionNameClaimAddress = "claim_address"

func (d ClaimAddress) TableName() string {
	return CollectionNameClaimAddress
}

func (d ClaimAddress) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:        []string{"-address"},
		Unique:     true,
		Background: true,
	})
	ensureIndexes(d.TableName(), indexes)
}

func (d ClaimAddress) PkKvPair() map[string]interface{} {
	return bson.M{"address": d.Address}
}

func (d ClaimAddress) Exist(address string) (bool, error) {
	var info ClaimAddress
	query := bson.M{
		"address": address,
	}
	fn := func(c *qmgo.Collection) error {
		return c.Find(_ctx, query).One(&info)
	}

	err := ExecCollection(d.TableName(), fn)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return false, nil
		}
		return false, err
	}

	return info.Address == address, nil
}

func (d ClaimAddress) Save(info ClaimAddress) error {
	fn := func(c *qmgo.Collection) error {
		_, err := c.InsertOne(_ctx, info)
		if err != nil {
			return err
		}
		return nil
	}

	return ExecCollection(d.TableName(), fn)
}

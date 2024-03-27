package model

import (
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
)

// Asset is an object representing the database table.
type Asset struct {
	Name             string    `bson:"name"`
	Symbol           string    `bson:"symbol"`
	Denom            string    `bson:"denom"`
	DenomLpt         string    `bson:"denom_lpt"`
	Offline          bool      `bson:"offline"`
	Scale            int       `bson:"scale"`
	CoinID           string    `bson:"coin_id"`
	Platform         string    `bson:"platform"`
	Protocol         int       `bson:"protocol"`
	Icon             string    `bson:"icon"`
	Tips             string    `bson:"tips"`
	IbcTransferInfos []IbcInfo `bson:"ibc_transfer_infos"`
	CreateTime       int64     `bson:"create_time"`
	UpdateTime       int64     `bson:"update_time"`
}

type (
	IbcInfo struct {
		InPath InPath      `bson:"in_path"`
		Traces []TraceInfo `bson:"traces"`
	}
	TraceInfo struct {
		Port     string `bson:"port"`
		Channel  string `bson:"channel"`
		ChainId  string `bson:"chain_id"`
		Denom    string `bson:"denom"`
		Platform string `bson:"platform"`
	}
	InPath struct {
		Port    string `bson:"port"`
		Channel string `bson:"channel"`
	}
)

const (
	_ = iota
	Native
	HashLock
	Ibc
)

const CollectionName = "coins"

func (d Asset) TableName() string {
	return CollectionName
}

func (d Asset) EnsureIndexes() {
	var indexes []options.IndexModel
	indexes = append(indexes, options.IndexModel{
		Key:        []string{"-denom"},
		Unique:     true,
		Background: true,
	})
	ensureIndexes(d.TableName(), indexes)
}

func (d Asset) PkKvPair() map[string]interface{} {
	return bson.M{"denom": d.Denom}
}

func (d Asset) findAllWithSelector(selector bson.M) ([]Asset, error) {
	var assets []Asset
	query := bson.M{}
	fn := func(c *qmgo.Collection) error {
		return c.Find(_ctx, query).Select(selector).All(&assets)
	}

	err := ExecCollection(d.TableName(), fn)

	return assets, err
}

func (d Asset) QueryAllScale() ([]Asset, error) {
	return d.findAllWithSelector(bson.M{"denom": 1, "scale": 1})
}

func (d Asset) FindAll() ([]Asset, error) {
	return d.findAllWithSelector(bson.M{})
}

func (d Asset) FindCoinByDenom(denom string) (Asset, error) {
	var assets Asset
	selector := bson.M{}
	query := bson.M{
		"denom": denom,
	}
	fn := func(c *qmgo.Collection) error {
		return c.Find(_ctx, query).Select(selector).One(&assets)
	}

	err := ExecCollection(d.TableName(), fn)

	return assets, err
}

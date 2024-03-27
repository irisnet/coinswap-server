package model

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Notice struct {
	Content  string `bson:"content" json:"content"`
	CreateAt int64  `bson:"create_at" json:"-"`
	UpdateAt int64  `bson:"update_at" json:"-"`
}

const CollectionNameNotice = "notice"

func (d Notice) TableName() string {
	return CollectionNameNotice
}

func (d Notice) EnsureIndexes() {
	//todo
}

func (d Notice) PkKvPair() map[string]interface{} {
	return bson.M{"content": d.Content}
}

func (w Notice) FindLatestCreateAtOne() (notice Notice, err error) {
	err = NewCollection(w).Find(_ctx, bson.M{}).Sort("-create_at").One(&notice)
	return
}

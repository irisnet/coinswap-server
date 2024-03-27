package model

import (
	"context"
	"fmt"
	"github.com/irisnet/coinswap-server/config"
	"github.com/irisnet/coinswap-server/internal/app/pkg/logger"
	"github.com/irisnet/coinswap-server/internal/app/server/types"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
)

var (
	_ctx  = context.Background()
	_conf *config.Config
	_cli  *qmgo.Client
)

type (
	Docs interface {
		// collection name
		TableName() string
		// ensure indexes
		EnsureIndexes()
		// primary key pair(used to find a unique record)
		PkKvPair() map[string]interface{}
	}
)

var (
	Collections = []Docs{
		new(Farms),
		new(Asset),
		new(ClaimAddress),
		new(StatisticInfo),
		new(WhiteList),
		new(LpTokens),
	}
)

func GetConf() *config.Config {
	if _conf == nil {
		logger.Fatal("db.Init not work")
	}
	return _conf
}

func GetClient() *qmgo.Client {
	return _cli
}

func Init(conf *config.Config) {
	_conf = conf
	var maxPoolSize uint64 = 4096
	// PrimaryMode indicates that only a primary is considered for reading. This is the default mode.
	client, err := qmgo.NewClient(_ctx, &qmgo.Config{
		Uri:         conf.MongoDb.NodeUri,
		Database:    conf.MongoDb.Database,
		MaxPoolSize: &maxPoolSize,
	})
	if err != nil {
		logger.Fatal(fmt.Sprintf("connect mongo failed, uri: %s, err:%s", conf.MongoDb.NodeUri, err.Error()))
	}
	_cli = client

	logger.Info("init db success")

	// ensure table indexes
	ensureDocsIndexes()
	return
}

func Close() {
	logger.Info("release resource :mongoDb")
	if _cli != nil {
		_cli.Close(_ctx)
	}
}

func ensureIndexes(collectionName string, indexes []options.IndexModel) {
	c := _cli.Database(GetConf().MongoDb.Database).Collection(collectionName)
	if len(indexes) > 0 {
		for _, v := range indexes {
			if err := c.CreateOneIndex(context.Background(), v); err != nil {
				logger.Warn("ensure index fail", logger.String("collectionName", collectionName),
					logger.String("index", types.MarshalJsonIgnoreErr(v)),
					logger.String("err", err.Error()))
			}
		}
	}
}

// get collection object
func ExecCollection(collectionName string, s func(*qmgo.Collection) error) error {
	c := _cli.Database(GetConf().MongoDb.Database).Collection(collectionName)
	return s(c)
}

func NewCollection(docs Docs) *qmgo.Collection {
	return _cli.Database(GetConf().MongoDb.Database).Collection(docs.TableName())
}

func ensureDocsIndexes() {
	if len(Collections) > 0 {
		for _, v := range Collections {
			v.EnsureIndexes()
		}
	}
}

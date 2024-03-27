module github.com/irisnet/coinswap-server

go 1.16

require (
	github.com/bsm/redislock v0.7.0
	github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5
	github.com/go-kit/kit v0.10.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/go-redis/redis/v8 v8.8.0
	github.com/gorilla/mux v1.8.0
	github.com/irisnet/irishub-sdk-go v0.0.0-20210902041149-3d7b782962c6
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/qiniu/qmgo v1.0.4
	github.com/robfig/cron/v3 v3.0.1
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/tendermint/tendermint v0.34.11
	github.com/volatiletech/sqlboiler/v4 v4.8.6
	go.mongodb.org/mongo-driver v1.7.2
	go.uber.org/zap v1.17.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

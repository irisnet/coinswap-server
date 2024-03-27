package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var conf Config

type (
	// Config define a struct for starting the http server
	Config struct {
		MongoDb Mongodb `mapstructure:"mongodb"`
		Server  Server  `mapstructure:"server"`
		Redis   Redis   `mapstructure:"redis"`
		IRIShub IRIShub `mapstructure:"irishub"`
		Task    Task    `mapstructure:"task"`
	}

	Mongodb struct {
		NodeUri  string `mapstructure:"node_uri"`
		Database string `mapstructure:"database"`
	}

	// Redis define a struct for redis server
	Redis struct {
		Address  string `mapstructure:"address"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}

	// Server define a struct for http server
	Server struct {
		Address              string `mapstructure:"address"`
		PriceDenom           string `mapstructure:"price_denom"`
		HandleFarmsWorkerNum int    `mapstructure:"handle_farms_worker_num"`
	}

	IRIShub struct {
		BaseDenom    string `mapstructure:"base_denom"`
		RcpAddr      string `mapstructure:"rpc_address"`
		GrpcAddr     string `mapstructure:"grpc_address"`
		LcdAddr      string `mapstructure:"lcd_address"`
		ChainID      string `mapstructure:"chain_id"`
		Fee          string `mapstructure:"fee"`
		BlockPerYear uint   `mapstructure:"block_per_year"`
	}

	Task struct {
		Enable                        bool `mapstructure:"enable"`
		CronTimeUpdateTotalVolumeLock int  `mapstructure:"cron_time_update_total_volume_lock"`
		CronTimeUpdateLiquidityPool   int  `mapstructure:"cron_time_update_liquidity_pool"`
		CronTimeUpdateFarm            int  `mapstructure:"cron_time_update_farm"`
	}
)

func Get() Config {
	return conf
}

func Load(cmd *cobra.Command, home string) error {
	rootViper := viper.New()
	_ = rootViper.BindPFlags(cmd.Flags())
	// Find home directory.
	rootViper.AddConfigPath(rootViper.GetString(home))
	rootViper.SetConfigName("config")
	rootViper.SetConfigType("toml")

	// Find and read the config file
	if err := rootViper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return err
	}

	if err := rootViper.Unmarshal(&conf); err != nil {
		return err
	}

	return nil
}

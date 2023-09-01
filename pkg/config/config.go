package config

import "github.com/spf13/viper"

type Config struct {
	ListenAddress  string             `mapstructure:"listen_address"`
	Redis          string             `mapstructure:"redis"`
	QueueSize      int                `mapstructure:"queue_size"`
	WorkerPoolSize int                `mapstructure:"worker_pool_size"`
	UserQuota      map[uint]UserQuota `mapstructure:"user_quota"`
}

type UserQuota struct {
	RequstPerMinute int64 `mapstructure:"request_per_minute"`
	SizePerMonth    int64 `mapstructure:"size_per_month"`
}

func Load() *Config {
	viper.SetConfigFile("./config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}
	return &config
}

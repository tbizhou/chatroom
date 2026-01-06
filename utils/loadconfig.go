package utils

import "github.com/spf13/viper"

type Config struct {
	Redis RedisConfig `mapstructure:"redis"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
	DB       int    `mapstructure:"db"`
	Port     int    `mapstructure:"port"`
	PoolSize int    `mapstructure:"pool_size"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // 文件名 (不带扩展名)
	viper.SetConfigType("yaml")   // 文件类型
	viper.AddConfigPath(".")      // 查找路径 (当前目录)

	// 自动读取环境变量 (可选，非常实用)
	// 比如环境变量设置 APP_REDIS_PASSWORD=123 即可覆盖 yaml 中的值
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

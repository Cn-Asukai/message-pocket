package config

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/samber/lo"
	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	URL     string `mapstructure:"url"`
	Token   string `mapstructure:"token"`
	GroupID string `mapstructure:"NAPCAT_GROUP_ID"`
}

var (
	instance *Config
	once     sync.Once
)

// LoadConfig 从环境变量加载配置
func LoadConfig() (*Config, error) {
	var err error
	once.Do(func() {
		// 初始化viper，只使用环境变量
		v := viper.New()

		v.AutomaticEnv() // 自动绑定环境变量

		// 将配置绑定到结构体
		instance = &Config{}
		if err = v.Unmarshal(instance); err != nil {
			return
		}

		if isEmpty := lo.IsEmpty(*instance); isEmpty {
			err = errors.New("some config is empty")
			return
		}

		slog.Debug("配置加载成功: %+v", instance)
	})
	return instance, err
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	if instance == nil {
		cfg, err := LoadConfig()
		if err != nil {
			panic(err)
		}
		return cfg
	}
	return instance
}

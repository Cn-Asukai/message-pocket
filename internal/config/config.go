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
	ServerConfig ServerConfig `yaml:"server" mapstructure:"server"`
	NapCatConfig NapCatConfig `yaml:"napcat" mapstructure:"napcat"`
}

type ServerConfig struct {
	OpenToken string `yaml:"open_token" mapstructure:"open_token"`
}

type NapCatConfig struct {
	URL     string `yaml:"url" mapstructure:"url"`
	Token   string `yaml:"token" mapstructure:"token"`
	GroupID string `yaml:"group_id" mapstructure:"group_id"`
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
		v.AddConfigPath("./")
		v.SetConfigName("config")
		v.SetConfigType("yaml")

		// 将配置绑定到结构体
		instance = &Config{}
		if err = v.ReadInConfig(); err != nil {
			return
		}
		if err = v.Unmarshal(instance); err != nil {
			return
		}

		if isEmpty := lo.IsEmpty(*instance); isEmpty {
			err = errors.New("some config is empty")
			return
		}

		slog.Debug("config loaded success.", instance)
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

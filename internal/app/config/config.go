package config

import (
	"log"

	"github.com/spf13/viper"
)

// 定义配置结构体
type Config struct {
	Service  ServiceConfig  `yaml:"service" mapstructure:"service"`
	Database DatabaseConfig `yaml:"database" mapstructure:"database"`
}

type ServiceConfig struct {
	Port string `yaml:"port" mapstructure:"port"`
}

type DatabaseConfig struct {
	Link        string `yaml:"link" mapstructure:"link"`
	MaxIdle     int    `yaml:"maxIdle" mapstructure:"maxIdle"`
	MaxOpen     int    `yaml:"maxOpen" mapstructure:"maxOpen"`
	MaxLifeTime int    `yaml:"maxOpen" mapstructure:"maxOpen"`
}

// 初始化 viper 并加载配置文件
func SetupConfig() (*Config, error) {
	viper.SetConfigName("config") // 配置文件名（不带扩展名）
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath("config") // 确保路径与 config.yaml 文件所在目录一致

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("加载配置文件失败：%v", err)
		return nil, err
	}

	// 将配置绑定到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("解析配置文件失败：%v", err)
		return nil, err
	}
	log.Println("初始化配置成功")
	return &config, nil
}

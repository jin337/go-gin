package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

// 设置默认值
var defaultConfig = map[string]interface{}{
	"Service.port":           "8080",
	"Service.TokenSecret":    "123456",
	"Database.maxIdle":       10,
	"Database.maxOpen":       100,
	"Database.maxLifeTime":   30,
	"Database.migrateTables": true,
	"Log.dirName":            "log",
}

// 定义配置结构体
type Config struct {
	Service  ServiceConfig  `yaml:"Service" mapstructure:"Service"`
	Database DatabaseConfig `yaml:"Database" mapstructure:"Database"`
	Log      LogConfig      `yaml:"Log" mapstructure:"Log"`
}

// 服务
type ServiceConfig struct {
	Port        string `yaml:"port" mapstructure:"port"`
	TokenSecret string `yaml:"tokenSecret" mapstructure:"tokenSecret"`
}

// 数据库
type DatabaseConfig struct {
	Link          string `yaml:"link" mapstructure:"link"`
	MaxIdle       int    `yaml:"maxIdle" mapstructure:"maxIdle"`
	MaxOpen       int    `yaml:"maxOpen" mapstructure:"maxOpen"`
	MaxLifeTime   int    `yaml:"maxLifeTime" mapstructure:"maxLifeTime"`
	MigrateTables bool   `yaml:"migrateTables" mapstructure:"migrateTables"`
}

// 日志
type LogConfig struct {
	DirName string `yaml:"dirName" mapstructure:"dirName"`
}

var (
	globalCfgOnce sync.Once
	globalCfg     *Config
)

// 初始化 viper 并加载配置文件
func SetupConfig(env string) error {
	log.Printf("当前环境：%v", env)
	// 配置文件名
	switch env {
	case "dev":
		viper.SetConfigName("config.dev")
	case "pro":
		viper.SetConfigName("config.pro")
	default:
		viper.SetConfigName("config.dev")
	}

	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath("config") // 路径

	// 设置默认值
	for k, v := range defaultConfig {
		viper.SetDefault(k, v)
	}

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("加载配置文件失败：%v", err)
		return err
	}

	// 将配置绑定到结构体
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Printf("解析配置文件失败：%v", err)
		return err
	}

	// 设置全局配置
	globalCfgOnce.Do(func() {
		globalCfg = &config
	})

	log.Println("初始化配置成功")
	return nil
}

// 获取全局配置
func GetGlobalConfig() *Config {
	return globalCfg
}

package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gorm.io/gorm/logger"
)

// ApplicationConfigEnv 环境应用配置
type ApplicationConfigEnv struct {

	// mysql 配置
	Mysql Mysql
	// logger 配置
	Logger Logger
	// identified 配置
	Snowflake Snowflake
}

// Mysql 配置
type Mysql struct {
	DSN         string
	MaxOpen     int
	MaxLifetime int
}

// Logger 配置
type Logger struct {
	Level logger.LogLevel
}

// Snowflake 配置
type Snowflake struct {
	StartTime int64
	NodeId    int64
}

// ApplicationConfig 应用配置
var ApplicationConfig *ApplicationConfigEnv

// 调用当前文件时，在变量和常量处理之后执行init方法
func init() {
	initialization()
	err := viper.Unmarshal(&ApplicationConfig)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct, %w", err))
	}
}

// 初始化配置（为了区分init方法定义为initialization）
func initialization() {
	// 激活环境
	var active string
	// 从命令行参数获取激活环境
	flag.StringVar(&active, "active", "dev", "The active environment")
	//	执行命令行参数解析
	flag.Parse()
	if active == "" {
		fmt.Println("Active from flag is \"\"")
		// 默认环境为dev
		active = "dev"
	}
	// 设置配置文件类型
	viper.SetConfigType("yaml")
	// 设置配置文件所在目录
	viper.AddConfigPath("./config")
	// 设置配置文件名称
	viper.SetConfigName("application_" + active)
	err := viper.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		ok := errors.As(err, &configFileNotFoundError)
		if ok {
			panic(fmt.Errorf("config file not found: %w", err))
		}
		panic(fmt.Errorf("config file was found but another error was produced: %w", err))
	}
	// 配置文件变化时，做出记录
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("config file change: ", in.Name)
		err := viper.Unmarshal(&ApplicationConfig)
		if err != nil {
			fmt.Println("unable to decode into struct")
		}
	})
	// 监控配置文件，当配置文件变化时应用新的配置文件
	viper.WatchConfig()
}

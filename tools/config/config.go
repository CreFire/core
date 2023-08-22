package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

func IsDev() bool {
	return Conf.Model == "development"
}

var Conf *Config

// config init
func init() {
	Conf = &Config{}
	viper.SetConfigName("config") // 设置配置文件名
	viper.AddConfigPath("./conf") // 设置配置文件路径
	err := viper.ReadInConfig()   // 读取配置文件
	if err != nil {
		log.Fatalf("Failed to read the config file: %s", err)
	}
	err = viper.Unmarshal(&Conf) // 解析配置文件
	if err != nil {
		log.Fatal("viper ReadConfig err", err, Conf)
		return
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Conf file changed:", e.Name)
	})
	viper.WatchConfig()
}

package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Name string `mapstructure:"name"`
}

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func main() {
	debug := GetEnvInfo("debug")
	v := viper.New()
	fileName := "./viper_test/config/config.yaml"
	if debug {
		fileName = "./viper_test/test/config.yaml"
	}
	v.SetConfigFile(fileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	serverConfig := ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	fmt.Println(serverConfig.Name)

	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置信息发生变化：", in.Name)
		err := v.ReadInConfig()
		if err != nil {
			panic(err)
		}
		if err := v.Unmarshal(&serverConfig); err != nil {
			panic(err)
		}
	})
}

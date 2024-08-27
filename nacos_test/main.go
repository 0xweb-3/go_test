package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type UserSrv struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JWTInfo struct {
	SigningKey string `mapstructure:"signing_key" json:"signing_key"`
}

type RedisInfo struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type Config struct {
	Name    string       `mapstructure:"name" json:"name"`
	Port    int          `mapstructure:"port" json:"port"`
	Host    string       `mapstructure:"host" json:"host"`
	UserSrv UserSrv      `mapstructure:"user_srv" json:"user_srv"`
	Jwt     JWTInfo      `mapstructure:"jwt" json:"jwt"`
	Redis   RedisInfo    `mapstructure:"redis" json:"redis"`
	Consul  ConsulConfig `mapstructure:"consul" json:"consul"`
}

func main() {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "192.168.21.6",
			Port:   8848,
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         "8e3fbfdb-41ed-4c52-936e-b8cf975e2d35", // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: false,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建动态配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user_web.json",
		Group:  "dev",
	})
	if err != nil {
		panic(err)
	}

	cnf := &Config{}
	json.Unmarshal([]byte(content), &cnf)
	fmt.Println(cnf.Jwt.SigningKey)

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: "user_web.json",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("配置文件发生变化")
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})

	if err != nil {
		panic(err)
	}

	select {} // 永久阻塞
}

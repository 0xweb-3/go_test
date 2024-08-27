package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

// 服务的注册
func Register(address, name, id string, port int, tags []string) {
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	err = client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Port:    port,
		Address: address,
		Check: &api.AgentServiceCheck{
			HTTP:                           "http://192.168.21.4:8081/helth", // 这里不能用127.0.0.1
			Timeout:                        "5s",
			Interval:                       "5s",
			DeregisterCriticalServiceAfter: "20s",
		},
	})
	if err != nil {
		panic(err)
	}
}

// 服务发现
func AllService() {
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 会获取所有服务
	srvs, err := client.Agent().Services()
	if err != nil {
		panic(err)
	}
	for key, srv := range srvs {
		fmt.Println(key, srv)
	}
	userWebsrv, err := client.Agent().ServicesWithFilter(`Service == "xin_user_web"`) // 这里使用的是name
	if err != nil {
		panic(err)
	}
	for key, srv := range userWebsrv {
		fmt.Println(key, srv)
	}
	//fmt.Println(fmt.Sprintf("%v", userWebsrv))
}

// 服务的注销
func Deregister(id string) {
	cfg := api.DefaultConfig()
	cfg.Address = "127.0.0.1:8500"

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	err = client.Agent().ServiceDeregister(id)
	if err != nil {
		panic(err)
	}
	log.Printf("Service %s deregistered", id)
}

func main() {
	//Register("127.0.0.1", "xin_user_web", "xin_user_web", 8081, []string{"xin", "user"})
	//AllService()
	Deregister("xin_user_web")
}

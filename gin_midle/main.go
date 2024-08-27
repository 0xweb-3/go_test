package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r := gin.Default()

	go func() {
		r.Run(":8081")
	}()
	// 如果想要接受到信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 后续处理逻辑
	fmt.Println("服务终止")
	// 做一些断开链接及清理处理
	fmt.Println("注销服务。。。。")
}

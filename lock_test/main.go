package main

import (
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"time"
)

func lock(index int) error {
	// 链接池
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "127.0.0.1:6380",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)

	mutexname := "my-global-mutex" // 锁名称
	mutex := rs.NewMutex(mutexname)

	// 上锁
	if err := mutex.Lock(); err != nil {
		return err
	}
	fmt.Println("已经锁定", index)

	// todo 实现需要进行锁操作的代码
	time.Sleep(time.Second)

	// 释放锁，以便其他进程或线程可以获得锁。
	if ok, err := mutex.Unlock(); !ok || err != nil {
		return err
	}
	fmt.Println("已经释放", index)

	return nil
}

func main() {
	var g errgroup.Group
	for i := 0; i < 20; i++ {
		i := i
		g.Go(func() error {
			fmt.Println("in index", i)
			return lock(i)
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Println("Error:", err)
	}
}

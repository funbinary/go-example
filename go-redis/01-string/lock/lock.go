package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type Lock struct {
	ctx context.Context
	rdb *redis.Client
	key string
}

func New(key string) *Lock {
	lock := &Lock{
		ctx: context.Background(),
		rdb: redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   0,
		}),
		key: key,
	}

	return lock
}

//
func (self *Lock) Acquire() bool {
	ok, _ := self.rdb.SetNX(self.ctx, self.key, "lock", -1).Result()
	return ok
}

func (self *Lock) Release() bool {
	res, err := self.rdb.Del(self.ctx, self.key).Result()
	return err == nil && res == 1
}

func main() {
	l := New("lock")
	if l.Acquire() {
		fmt.Println("获取成功")
	} else {
		fmt.Println("获取失败")
	}
	if l.Acquire() {
		fmt.Println("获取成功")
	} else {
		fmt.Println("获取失败")
	}
	if l.Release() {
		fmt.Println("释放成功")
	} else {
		fmt.Println("释放失败")
	}
	if l.Acquire() {
		fmt.Println("获取成功")
	} else {
		fmt.Println("获取失败")
	}
}

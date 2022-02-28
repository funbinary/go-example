package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func main() {
	var err error
	ctx := context.Background()
	//连接方式1
	//rdb := redis.NewClient(&redis.Options{
	//	Addr: "127.0.0.1:6379",
	//	DB:   0,
	//})

	//连接方式2
	//使用字符串连接  "redis://<user>:<pass>@127.0.0.1:6379/<db>"
	opt, err := redis.ParseURL("redis://@127.0.0.1:6379/0")
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	//fmt.Println(pong)

	//设置Key
	err = rdb.Set(ctx, "name", "redis", 1*time.Second).Err()
	if err != nil {
		fmt.Println(err)
	}

	//获取过期时间
	tm, err := rdb.TTL(ctx, "name").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tm)

	//获取key
	name, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(name)

	//判断key是否存在
	val, err := rdb.Get(ctx, "None").Result()
	if err == redis.Nil {
		fmt.Println("None 不存在")
	} else if err != nil {
		fmt.Println(val, err)
	} else {
		fmt.Println(val)
	}

	//当key不存在时设置
	ok, err := rdb.SetNX(ctx, "count", 0, 1*time.Second).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ok)
	}

	//对值加1
	result, err := rdb.Incr(ctx, "count").Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	//对值加上value
	result, err = rdb.IncrBy(ctx, "count", 5).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

}

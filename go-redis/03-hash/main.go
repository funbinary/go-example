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
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	//fmt.Println(pong)

	//添加
	doupo := map[string]interface{}{
		"title":     "斗破苍穹",
		"content":   "林动乱入...",
		"author":    "天蚕土豆",
		"create_at": time.Now(),
	}
	if err = rdb.HMSet(ctx, "article:10086", doupo).Err(); err != nil {
		fmt.Println(err)
	}

	//获取标题和内容
	fmt.Println(rdb.HMGet(ctx, "article:10086", "title", "content").Result())

	//获取所有内容
	fmt.Println(rdb.HGetAll(ctx, "article:10086").Result())

	//判断字段是否存在
	fmt.Println(rdb.HExists(ctx, "article:10086", "about").Result())
	fmt.Println(rdb.HExists(ctx, "article:10086", "title").Result())

	//如果key不存在则设置
	fmt.Println("如果key不存在则设置")
	fmt.Println(rdb.HSetNX(ctx, "article:10086", "title", "武动乾坤").Result())
	fmt.Println(rdb.HSetNX(ctx, "article:10086", "id", 10086).Result())

	// 删除
	fmt.Println("删除")
	fmt.Println(rdb.HDel(ctx, "article:10086", "id").Result())
	fmt.Println(rdb.HDel(ctx, "article:10086", "none").Result())
}

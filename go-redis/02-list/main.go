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

	//向list尾部添加值
	err = rdb.RPush(ctx, "books", "Redis设计与实现").Err()
	if err != nil {
		fmt.Println(err)
	}
	//向list尾部添加值
	err = rdb.RPush(ctx, "books", "Redis实战").Err()
	if err != nil {
		fmt.Println(err)
	}
	//为指定索引添加新元素,如果索引无效则失败
	err = rdb.LSet(ctx, "books", 2, "无效的索引").Err()
	if err != nil {
		fmt.Println(err) //索引无效失败
	}

	err = rdb.LSet(ctx, "books", 1, "Redis使用手册").Err()
	if err != nil {
		fmt.Println(err) //索引有效成功
	}

	// 如果count=0，将移除列表中所有包含指定value的元素
	// 如果count大于0，将从列表头部开始移除最先发现的count个元素
	// 如果count小于0，将从列表尾部开始移除最先发现的count(绝对值)个元素
	// 返回包含移除元素的数量
	fmt.Println(rdb.LRem(ctx, "books", 1, "Redis实战").Result())

	//获取列表books的长度
	booksLen, err := rdb.LLen(ctx, "books").Result()
	fmt.Println(booksLen, err)

	//遍历列表books
	list, err := rdb.LRange(ctx, "books", 0, booksLen-1).Result()
	fmt.Println(list, err)

	//阻塞式头部弹出操作
	fmt.Println(rdb.BLPop(ctx, 1*time.Second, "books").Result())

}

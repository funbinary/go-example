package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io"
	"os"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	ctx := context.Background()
	rdb.FlushDB(ctx)
	file, err := os.Open("D://workspace//myOpenSource//go-example//go-redis//01-string//cache//redis-logo.jpg")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	rdb.Set(ctx, "redis-logo.jpg", data, -1)
	fmt.Println(rdb.Get(ctx, "redis-logo.jpg"))
}

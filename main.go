package main

import (
	"github.com/bin-work/go-example/pkg/bfile"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath(bfile.SelfDir())
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

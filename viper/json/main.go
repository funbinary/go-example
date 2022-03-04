package main

import (
	"fmt"
	"github.com/bin-work/go-example/internal/file"
	"github.com/spf13/viper"
)

type Host struct {
	Address string `mapstructure:address`
	Port    int    `mapstructure:port`
}

type Metric struct {
	Host string `mapstructure:host`
	Port int    `mapstructure:port`
}

type Warehouse struct {
	Host string `mapstructure:host`
	Port int    `mapstructure:port`
}

type Datastore struct {
	Metric    Metric
	Warehouse Warehouse
}

type Config struct {
	Host      Host
	Datastore Datastore
}

func main() {
	viper.SetConfigName("config")
	fmt.Println(file.ExtName("config.json"))
	viper.SetConfigType("json")
	viper.AddConfigPath(file.SelfDir())
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.GetString("datastore.metric.host"))
	var config Config
	err = viper.Unmarshal(&config)
	viper.Set("mysql_ip", "222")
	viper.WriteConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.IsSet("datastore")) //true
	viper.MergeConfigMap(map[string]interface{}{
		"test": 1,
	})
	viper.Debug()
	//viper.Reset()
	//fmt.Println("=========reset==================")
	//viper.Debug()

	viper.SafeWriteConfig()
}

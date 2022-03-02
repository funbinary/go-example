package main

import (
	"fmt"
	"github.com/bin-work/go-example/internal/file"
	"github.com/spf13/viper"
	"os"
)

type Mysql struct {
	IP       string `mapstructure:ip`
	Port     int    `mapstructure:port`
	User     string
	Password string
	Database string
}

type Config struct {
	Mysql Mysql
}

type myFlag struct{}

func (f myFlag) HasChanged() bool    { return false }
func (f myFlag) Name() string        { return "my-flag-name" }
func (f myFlag) ValueString() string { return "my-flag-value" }
func (f myFlag) ValueType() string   { return "string" }

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath(file.SelfDir())
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}
	//fmt.Println(viper.AllKeys())
	//fmt.Println(viper.AllSettings())
	fmt.Println(viper.GetString("mysql.ip"))
	// 绑定Flag
	viper.BindFlagValue("mysql.ip", myFlag{})
	fmt.Println("BindFlagValue:", viper.GetString("mysql.ip"))
	// 绑定环境变量
	viper.BindEnv("mysql.ip")
	val, ok := os.LookupEnv("mysql.ip")

	fmt.Println(val, ok)
	fmt.Println(viper.GetString("mysql.ip"))
	//fmt.Println(config.Mysql.IP)
	//fmt.Println(config.Mysql.Port)
	//fmt.Println(config.Mysql.User)
	//fmt.Println(config.Mysql.Password)
	//fmt.Println(config.Mysql.Database)
}

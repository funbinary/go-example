package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Account struct {
	gorm.Model
	HeadImg     string
	Username    string
	DisplayName string
	Sex         string
	NickName    string
	Email       string
}

func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	log.Println("------------------")
	db.AutoMigrate(&Account{})
	log.Println("------------------")
	db.Create(&Account{

		HeadImg:     "",
		Username:    "h",
		DisplayName: "hh",
	})

}

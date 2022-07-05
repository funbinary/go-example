package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	MediaServerId string `json:"mediaServerId,omitempty"`
	App           string `json:"app,omitempty"`
	FileName      string `json:"fileName,omitempty"`
	FilePath      string `json:"filePath,omitempty"`
	FileSize      string `json:"fileSize,omitempty"`
	Folder        string `json:"folder,omitempty"`
	StartTime     int    `json:"startTime,omitempty"`
	Stream        string `json:"stream,omitempty"`
	TimeLen       int    `json:"timeLen,omitempty"`
	Url           string `json:"url,omitempty"`
	Vhost         string `json:"vhost,omitempty"`
}

func (self Resp) String() string {
	s := "Id:" + self.MediaServerId + "\n"
	s += "App:" + self.App + "\n"
	s += "FileName:" + self.FileName + "\n"
	s += "FilePath:" + self.FilePath + "\n"
	s += "FileSize:" + self.FileSize + "\n"
	s += "Folder:" + self.Folder + "\n"
	s += "StartTime:" + strconv.Itoa(self.StartTime) + "\n"
	s += "Stream:" + self.Stream + "\n"
	s += "TimeLen:" + strconv.Itoa(self.TimeLen) + "\n"
	s += "url:" + self.Url + "\n"
	s += "vhost:" + self.Vhost + "\n"
	return s
}

func main() {
	r := gin.Default()
	r.POST("/on_record_mp4", func(c *gin.Context) {
		resp := &Resp{}
		c.ShouldBindJSON(resp)
		fmt.Println(resp)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
		})
	})
	r.Run("192.168.3.100:8095") // 监听并在 0.0.0.0:8080 上启动服务
}

package main

import (
	"log"

	"github.com/XieXianbin/sms-provider/smtp"
)

func main() {
	s := smtp.New("smtp.qq.com:25", "xx@qq.com", "xx")
	log.Println(s.SendMail("xx@qq.com",
		"xx@163.com;",
		"这是subject",
		"这是body,<font color=red>red</font>"))
}

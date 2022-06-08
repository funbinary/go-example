package main

import (
	"github.com/bin-work/go-example/pkg/bfile"
	"github.com/bin-work/go-example/pkg/blogrus"

	"github.com/sirupsen/logrus"
)

func main() {
	writer := blogrus.NewDefaultMultiWriter(bfile.Join(bfile.SelfDir(), "/example.log"), 1<<26, 1000000)
	if writer == nil {
		panic("初始化日志文件失败")
	}
	logrus.SetOutput(writer)

	logrus.SetLevel(logrus.DebugLevel)

	format := blogrus.NewDefaultBFormatter()
	format.ShowFullLevel = false
	logrus.SetReportCaller(true)
	logrus.SetFormatter(format)
	logrus.Debugf("this is debug")
	logrus.Infof("This is info")
	func() {
		logrus.Infof("This is subFunc")
	}()
	logrus.Warn("This is warn")
	logrus.Warning("This is warning")
	logrus.Errorf("This is error")
	logrus.Fatal("This is Fatal")

}

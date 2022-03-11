/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
	"log"
	"os"
)

var (
	Config string //配置文件路径
	Path   string //FTP共享路径
)
var rootCmd = &cobra.Command{
	Use:   "ftp-server",
	Short: "快捷启动一个FTP服务器",
	Long:  `快速运行一个FTP服务器,For example:`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		dirver, err := file.NewDriver(Path)
		if err != nil {
			log.Fatal(err)
		}
		s, err := server.NewServer(&server.Options{
			Driver: dirver,
			Auth: &server.SimpleAuth{
				Name:     "admin",
				Password: "admin",
			},
			Port: 21,
			Perm: server.NewSimplePerm("root", "root"),
		})
		if err != nil {
			log.Fatal(err)
		}
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&Config, "config", "c", "config", "配置文件名称，配置文件需要与当前程序同一个目录")
	rootCmd.PersistentFlags().StringVarP(&Path, "path", "p", "./", "FTP文件目录")
	//rootCmd.PersistentFlags().StringVarP(&Path, "path", "p", "./", "FTP文件目录")

}

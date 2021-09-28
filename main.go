package main

import (
	"singo/api"
	"singo/conf"
	"singo/server"

	"github.com/robfig/cron"
)

func main() {
	// 从配置文件读取配置
	conf.Init()

	c := cron.New()
	_ = c.AddFunc("0 0 15 * * *", func() {
		api.FwAlert()
	})
	_ = c.AddFunc("0 0 10 * * *", func() {
		api.FwAlert()
	})
	//_ = c.AddFunc("5 * * * * *", func() {
	//	api.FwAlert()
	//})
	c.Start()

	// 装载路由
	r := server.NewRouter()
	_ = r.Run(":3000")
}

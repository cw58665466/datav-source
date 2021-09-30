package server

import (
	"io"
	"log"
	"os"
	"singo/api"
	"singo/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	//r.Use(gin.Logger())
	// 用户登录
	r.POST("/api/user/login", api.UserLogin)
	//r.GET("/api/user/home", middleware.JWTAuthMiddleware(), api.HomeHandler)

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	// 中间件, 顺序不能改
	//r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	//r.Use(middleware.CurrentUser())
	r.GET("/long_async", func(c *gin.Context) {
		// 创建在 goroutine 中使用的副本
		cCp := c.Copy()
		go func() {
			// 用 time.Sleep() 模拟一个长任务。
			time.Sleep(15 * time.Second)

			// 请注意您使用的是复制的上下文 "cCp"，这一点很重要
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})
	r.GET("/long_async_cp", func(c *gin.Context) {
		// 创建在 goroutine 中使用的副本
		cCp := c.Copy()
		go func() {
			// 用 time.Sleep() 模拟一个长任务。
			time.Sleep(15 * time.Second)

			// 请注意您使用的是复制的上下文 "cCp"，这一点很重要
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})

	r.GET("/long_sync", func(c *gin.Context) {
		// 用 time.Sleep() 模拟一个长任务。
		time.Sleep(5 * time.Second)

		// 因为没有使用 goroutine，不需要拷贝上下文
		log.Println("Done! in path " + c.Request.URL.Path)
	})
	r.GET("licenseAdd/:platform", api.LicenseAdd)
	//r.GET("fwAlert", api.FwAlert)
	// 路由
	v1 := r.Group("/api/datav")
	{
		v1.GET("amtArea/:type", api.AmtArea)
		v1.GET("countAmtArea/:type", api.CountAmtArea)
		v1.GET("amtAll/:type", api.AmtAll)
		v1.GET("amtTwo/:type", api.AmtTwo)
		v1.GET("amtCompare/:type", api.AmtCompare)
		v1.GET("amtTwoCompare/:type", api.AmtTwoCompare)
		v1.GET("twoCompare/:type", api.TwoCompare)
		v1.GET("duplicateCompanyUser/:type", api.DuplicateCompanyUser)
		v1.GET("attList/:type/:limit", api.AttList)
		v1.GET("companyUpUser/:type/:limit", api.CompanyUpUser)
		v1.GET("sumPrice/:type/:limit", api.SumPrice)
		v1.GET("companyRoom/:type/:startRoomCount/:endRoomCount/:limit", api.CompanyRoom)

		// 用户登录
		v1.POST("user/register", api.UserRegister)

		// 需要登录保护的
		auth := r.Group("/api")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			// User Routing
			auth.GET("user/me", api.HomeHandler)
		}
	}
	return r
}

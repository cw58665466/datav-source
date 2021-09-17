package server

import (
	"os"
	"singo/api"
	"singo/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 中间件, 顺序不能改
	r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	r.Use(middleware.CurrentUser())

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
		v1.GET("attList/:type", api.AttList)

		// 用户登录
		v1.POST("user/register", api.UserRegister)

		// 用户登录
		v1.POST("user/login", api.UserLogin)

		// 需要登录保护的
		auth := v1.Group("")
		auth.Use(middleware.AuthRequired())
		{
			// User Routing
			auth.GET("user/me", api.UserMe)
			auth.DELETE("user/logout", api.UserLogout)
		}
	}
	return r
}

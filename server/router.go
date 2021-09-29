package server

import (
	"singo/api"
	"singo/middleware"

	"github.com/gin-gonic/gin"
)

// NewRouter 路由配置
func NewRouter() *gin.Engine {
	r := gin.Default()

	// 用户登录
	r.POST("/api/user/login", api.UserLogin)
	//r.GET("/api/user/home", middleware.JWTAuthMiddleware(), api.HomeHandler)

	// 中间件, 顺序不能改
	//r.Use(middleware.Session(os.Getenv("SESSION_SECRET")))
	r.Use(middleware.Cors())
	//r.Use(middleware.CurrentUser())

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

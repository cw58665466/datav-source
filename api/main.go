package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v8"
	"singo/cache"
	"singo/conf"
	"singo/model"
	"singo/serializer"
	"singo/util"
	"time"
)
type Result struct {
	CompanyName string
}
// Ping 状态检查页面
func Ping(c *gin.Context) {
	c.JSON(200, serializer.Response{
		Code: 0,
		Msg:  "Pong",
	})
}
func Query(c *gin.Context)  {
	var logger = util.Log()
	var redisKey = "cache1"
	var db = model.DB
	var areaPrice []AreaPrice
	var rd = cache.RedisClient
	var ext = rd.Exists(redisKey)
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		logger.Info(jsonStr.Val())
		_ = json.Unmarshal([]byte(jsonStr.Val()), &areaPrice)

	} else {

		db.Raw("select left(c.parent_key, 6) as area_code,sum(a.sumPrice) as sum_price,c.area_name from (  select SUM(price)/100 as sumPrice,company_id from x_user_ticket where used_time >= \"2021-08-30 12:00:00\" and used_time <= \"2021-09-01 12:00:00\" and status = 1 and source_type =1 and price > 0  group by company_id ) a left join x_api_company b on a.company_id = b.id  left join x_police_station c on c.`key` = b.police_station_key group by c.parent_key").Scan(&areaPrice)

		data, _ := json.Marshal(areaPrice)
		rd.Set(redisKey,data,120*time.Second)


	}


	c.JSON(200, areaPrice)
}

// CurrentUser 获取当前用户
func CurrentUser(c *gin.Context) *model.XPoliceUser {
	if user, _ := c.Get("user"); user != nil {
		if u, ok := user.(*model.XPoliceUser); ok {
			return u
		}
	}
	return nil
}

// ErrorResponse 返回错误消息
func ErrorResponse(err error) serializer.Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := conf.T(fmt.Sprintf("Field.%s", e.Field))
			tag := conf.T(fmt.Sprintf("Tag.Valid.%s", e.Tag))
			return serializer.ParamErr(
				fmt.Sprintf("%s%s", field, tag),
				err,
			)
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.ParamErr("JSON类型不匹配", err)
	}

	return serializer.ParamErr("参数错误", err)
}

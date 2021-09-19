/**
编辑器到期提醒
*/
package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"singo/cache"
	"time"
)

const platformKey = "rainbowBaby"

func FwAlert() {
	var rd = cache.RedisClient
	platformList := rd.LRange(platformKey, 0, 10)
	listArr, _ := platformList.Result()
	for _, value := range listArr {
		redisKey := "platformTimeKey:" + value
		var ext = rd.Exists(redisKey)
		if ext.Val() == 1 {
			CheckPlatformLicense(value, rd.Get(redisKey).Val())
		} else {
			ErrorPush(rd.Get(redisKey).Val())
		}

	}

	//fmt.Println("test2")
}

func CheckPlatformLicense(platformName string, resetTime string) {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", resetTime, time.Local)
	noticeDate := t.AddDate(0, 0, 25)
	expiredDate := t.AddDate(0, 0, 30)
	if noticeDate.Before(time.Now()) {
		days := (int)(expiredDate.Unix()-time.Now().Unix()) / 3600 / 24
		url := fmt.Sprintf("https://fwalert.com/9db88a6c-9d97-4ea6-8b20-3f42b35ee3ab?platform=%s&days=%d", platformName, days)
		fmt.Println(Get(url))
	} else {
		fmt.Printf("%s -> %s", platformName, expiredDate)
	}

}
func NormalPush(platformName string, days int) {

}
func ErrorPush(platformName string) {

}
func Get(url string) string {
	res, err := http.Get(url)
	if err != nil {
		return ""
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return ""
	}
	return string(robots)
}
func LicenseAdd(c *gin.Context) {
	platform := c.Param("platform")
	platformTimeKey := "platformTimeKey:" + platform
	var rd = cache.RedisClient
	//
	//rd.Set(redisKey, data, redisCacheTime*time.Second)

	platformList := rd.LRange(platformKey, 0, 10)
	listArr, _ := platformList.Result()
	if !in(platform, listArr) {
		rd.LPush(platformKey, platform)
	}
	date := time.Now().Format("2006-01-02 15:04:05")
	rd.Set(platformTimeKey, date, 0)

	c.JSON(200, "保存成功")

}
func in(target string, strArray []string) bool {
	for _, element := range strArray {
		if target == element {
			return true
		}
	}
	return false
}

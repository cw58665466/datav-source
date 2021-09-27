package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"math"
	"singo/cache"
	"singo/model"
	"singo/util"
	"strconv"
	"time"
)

type CompanyResult struct {
	X string  `json:"x"`
	Y float64 `json:"y"`
}
type Chart struct {
	X string `json:"x"`
	Y int    `json:"y"`
	Z string `json:"z"`
}

type PriceSum struct {
	SumPrice  float64 `json:"sumPrice"`
	CompanyId int     `json:"company_id"`
	Nickname  string  `json:"nickname"`
}
type AmtResult struct {
	Value float64 `json:"value"`
}

type AttInfo struct {
	CompanyId   int    `json:"companyId"`
	CompanyName string `json:"companyName"`
	A           int    `json:"a"`
	C           int    `json:"c"`
	B           int    `json:"b"`
}

type AreaPrice struct {
	SumPrice     float64 `json:"sumPrice"`
	TwoSumPrice  float64 `json:"twoSumPrice"`
	P            float64 `json:"p"`
	AreaName     string  `json:"areaName"`
	AreaCode     string  `json:"areaCode"`
	CompanyCount int     `json:"companyCount"`
}

type CompanyRoomInfo struct {
	Id        int64   `json:"id"`
	AvgRoom   int64   `json:"avgRoom"`
	N         string  `json:"n"`
	AvgUpUser int64   `json:"avgUpUser"`
	RoomCount int64   `json:"roomCount"`
	Rjv       float64 `json:"rjv"`
}
type CompanyUpUserInfo struct {
	Nickname    string
	UpUserCount int64 `json:"upUserCount"`
}

type DuplicateCompanyUserInfo struct {
	Area string `json:"area"`
	Pv   int64  `json:"pv"`
}

//查询结果在redis中存储时间(秒)
const redisCacheTime = 600

func AmtArea(c *gin.Context) {
	var logger = util.Log()
	dateType := c.Param("type")
	var redisKey = "cache1-" + dateType
	var db = model.DB
	var areaPrice []AreaPrice
	var rd = cache.RedisClient
	var ext = rd.Exists(redisKey)
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		logger.Info(jsonStr.Val())
		_ = json.Unmarshal([]byte(jsonStr.Val()), &areaPrice)

	} else {

		db.Raw("SELECT\n\tsum(a.sumPrice) as sum_price,\n\tsum(d.sumPrice) as two_sum_price,\n\tsum(d.sumPrice)/sum(a.sumPrice) * 100 as p,\n\tleft(c.parent_key, 6) AS area_code,\n\tc.area_name\nFROM (\n\tSELECT\n\t\tSUM(price) / 100 AS sumPrice,\n\t\tcompany_id\n\tFROM\n\t\tx_user_ticket\n\tWHERE\n\t\tused_time >= DATE_SUB(NOW(), INTERVAL ? DAY)\n\t\tAND status = 1\n\t\tAND source_type = 1\n\t\tAND price > 0\n\tGROUP BY\n\t\tcompany_id) a\n\t\t\n\tleft join (SELECT\n\t\tSUM(price) / 100 AS sumPrice,\n\t\tcompany_id\n\tFROM\n\t\tx_user_ticket\n\tWHERE\n\t\tused_time >= DATE_SUB(NOW(), INTERVAL ? DAY)\n\t\tAND status = 1\n\t\tAND source_type = 1\n\t\tAND price > 0\n\t\tand company_id in (select id from x_api_company where company_service_amt = 0)\n\tGROUP BY\n\t\tcompany_id ) d on d.company_id = a.company_id\n\tLEFT JOIN x_api_company b ON a.company_id = b.id\n\tLEFT JOIN x_police_station c ON c. `key` = b.police_station_key\nGROUP BY\n\tc.parent_key", dateType, dateType).Scan(&areaPrice)

		for key, value := range areaPrice {
			if value.AreaCode == "330192" {
				areaPrice[key].AreaCode = "330114"
			}
			if value.AreaCode == "330183" {
				areaPrice[key].AreaCode = "330111"
			}
			if value.AreaCode == "330185" {
				areaPrice[key].AreaCode = "330112"
			}

		}

		data, _ := json.Marshal(areaPrice)
		rd.Set(redisKey, data, redisCacheTime*time.Second)

	}

	c.JSON(200, areaPrice)
}

func AmtAll(c *gin.Context) {
	//var logger = util.Log()
	dateType := c.Param("type")
	var redisKey = "cache2-" + dateType
	var db = model.DB
	var rd = cache.RedisClient
	var result AmtResult
	var ext = rd.Exists(redisKey)
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		result.Value, _ = strconv.ParseFloat(jsonStr.Val(), 64)

	} else {

		db.Raw("select SUM(price) / 100 AS value from x_user_ticket where used_time >= DATE_SUB(NOW(), INTERVAL ? DAY)", dateType).Scan(&result)

		rd.Set(redisKey, result.Value, redisCacheTime*time.Second)

	}
	var resultArr []AmtResult
	resultArr = append(resultArr, result)

	c.JSON(200, resultArr)
}
func AmtTwo(c *gin.Context) {
	//var logger = util.Log()
	dateType := c.Param("type")
	var redisKey = "cache3-" + dateType
	var db = model.DB
	var rd = cache.RedisClient
	var result AmtResult
	var ext = rd.Exists(redisKey)
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		result.Value, _ = strconv.ParseFloat(jsonStr.Val(), 64)

	} else {

		db.Raw("select  SUM(price) / 100 AS value from x_user_ticket where used_time >= DATE_SUB(NOW(), INTERVAL ? DAY) and company_id  in (select id from x_api_company where company_service_amt = 0)", dateType).Scan(&result)

		rd.Set(redisKey, result.Value, redisCacheTime*time.Second)

	}

	var resultArr []AmtResult
	resultArr = append(resultArr, result)

	c.JSON(200, resultArr)
}

func AmtCompare(c *gin.Context) {
	dateType := c.Param("type")
	var redisKey = "cache4-" + dateType
	var sum1, sum2 float64
	var rd = cache.RedisClient
	var result AmtResult
	intDateType, _ := strconv.Atoi(dateType)
	compareType := intDateType + intDateType

	var ext = rd.Exists(redisKey)

	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		result.Value, _ = strconv.ParseFloat(jsonStr.Val(), 64)
	} else {

		sumChan1 := make(chan float64, 1)
		sumChan2 := make(chan float64, 1)
		defer close(sumChan1)
		defer close(sumChan2)
		go SqlGet("select SUM(price) / 100 AS value from x_user_ticket where used_time >= DATE_SUB(NOW(), INTERVAL ? DAY)", dateType, sumChan1)
		go SqlGet("SELECT SUM(price) / 100 AS value FROM x_user_ticket WHERE used_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND used_time < DATE_SUB(NOW(), INTERVAL "+dateType+" DAY)", strconv.Itoa(compareType), sumChan2)
		sum1 = <-sumChan1
		sum2 = <-sumChan2
		//db.Raw("",dateType).Scan(&sum1)
		//db.Raw("SELECT SUM(price) / 100 AS value FROM x_user_ticket WHERE used_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND used_time < DATE_SUB(NOW(), INTERVAL ? DAY)",compareType , intDateType).Scan(&sum2)

		result.Value = (sum1 - sum2) / sum2 * 100

		rd.Set(redisKey, result.Value, redisCacheTime*time.Second)

	}

	var resultArr []AmtResult
	resultArr = append(resultArr, result)

	c.JSON(200, resultArr)

}
func TwoCompare(c *gin.Context) {
	dateType := c.Param("type")
	var redisKey = "cache5-" + dateType
	var re []CompanyResult
	//var sum1,sum2 float64
	var rd = cache.RedisClient
	var ext = rd.Exists(redisKey)
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		_ = json.Unmarshal([]byte(jsonStr.Val()), &re)

	} else {
		sum := make(chan float64, 2)
		defer close(sum)
		go SqlGet("select count(DISTINCT(company_id)) as value from x_user_ticket where used_time >= DATE_SUB(NOW(), INTERVAL ? DAY)", dateType, sum)
		go SqlGet("select count(DISTINCT(company_id)) as value from x_user_ticket where used_time >= DATE_SUB(NOW(), INTERVAL ? DAY)  and company_id  in (select id from x_api_company where company_service_amt = 0)", dateType, sum)

		sum1, sum2 := <-sum, <-sum
		var midSum float64
		if sum1 > sum2 {
			midSum = sum1
			sum1 = sum2
			sum2 = midSum
		}

		var two CompanyResult
		var noTwo CompanyResult
		two.X = "两元台票场所数"
		two.Y = sum1
		noTwo.X = "非两元台票场所数"
		noTwo.Y = sum2 - sum1
		re = append(re, two)
		re = append(re, noTwo)

		data, _ := json.Marshal(re)
		rd.Set(redisKey, data, redisCacheTime*time.Second)

	}

	c.JSON(200, re)

}
func AmtTwoCompare(c *gin.Context) {
	dateType := c.Param("type")
	var redisKey = "cache6-" + dateType
	var sum1, sum2 float64
	var rd = cache.RedisClient
	var result AmtResult
	intDateType, _ := strconv.Atoi(dateType)
	compareType := intDateType + intDateType

	var ext = rd.Exists(redisKey)

	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		result.Value, _ = strconv.ParseFloat(jsonStr.Val(), 64)
	} else {

		sumChan1 := make(chan float64, 1)
		sumChan2 := make(chan float64, 1)
		defer close(sumChan1)
		defer close(sumChan2)
		go SqlGet("select SUM(price) / 100 AS value from x_user_ticket where company_id  in (select id from x_api_company where company_service_amt = 0) and used_time >= DATE_SUB(NOW(), INTERVAL ? DAY)", dateType, sumChan1)
		go SqlGet("SELECT SUM(price) / 100 AS value FROM x_user_ticket WHERE company_id  in (select id from x_api_company where company_service_amt = 0) and used_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND used_time < DATE_SUB(NOW(), INTERVAL "+dateType+" DAY)", strconv.Itoa(compareType), sumChan2)
		sum1 = <-sumChan1
		sum2 = <-sumChan2
		//db.Raw("",dateType).Scan(&sum1)
		//db.Raw("SELECT SUM(price) / 100 AS value FROM x_user_ticket WHERE used_time >= DATE_SUB(NOW(), INTERVAL ? DAY) AND used_time < DATE_SUB(NOW(), INTERVAL ? DAY)",compareType , intDateType).Scan(&sum2)

		result.Value = (sum1 - sum2) / sum2 * 100

		rd.Set(redisKey, result.Value, redisCacheTime*time.Second)

	}

	var resultArr []AmtResult
	resultArr = append(resultArr, result)

	c.JSON(200, resultArr)

}
func CountAmtArea(c *gin.Context) {
	var logger = util.Log()
	dateType := c.Param("type")
	var redisKey = "cache7-" + dateType
	var db = model.DB
	var areaPrice []AreaPrice
	var rd = cache.RedisClient
	var ext = rd.Exists(redisKey)
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		logger.Info(jsonStr.Val())
		_ = json.Unmarshal([]byte(jsonStr.Val()), &areaPrice)

	} else {

		db.Raw("select left(c.parent_key, 6) as area_code,count(a.company_id) AS company_count,c.area_name from (  select SUM(price)/100 as sumPrice,company_id from x_user_ticket where used_time >= DATE_SUB(NOW(), INTERVAL ? DAY) and status = 1 and source_type =1 and price > 0  group by company_id ) a left join x_api_company b on a.company_id = b.id  left join x_police_station c on c.`key` = b.police_station_key group by c.parent_key", dateType).Scan(&areaPrice)

		for key, value := range areaPrice {
			if value.AreaCode == "330192" {
				areaPrice[key].AreaCode = "330114"
			}
			if value.AreaCode == "330183" {
				areaPrice[key].AreaCode = "330111"
			}
			if value.AreaCode == "330185" {
				areaPrice[key].AreaCode = "330112"
			}

		}

		data, _ := json.Marshal(areaPrice)
		rd.Set(redisKey, data, redisCacheTime*time.Second)

	}

	c.JSON(200, areaPrice)
}

func AttList(c *gin.Context) {

	//var logger = util.Log()
	dateType := c.Param("type")
	limit, _ := strconv.Atoi(c.Param("limit"))
	var redisKey = "cache8-" + dateType + "-" + c.Param("limit")
	var db = model.DB
	var attInfo []AttInfo
	var rd = cache.RedisClient
	var ext = rd.Exists(redisKey)

	var chart []Chart
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		//logger.Info(jsonStr.Val())
		_ = json.Unmarshal([]byte(jsonStr.Val()), &chart)

	} else {
		dateType := c.Param("type")
		intDateType, _ := strconv.Atoi(dateType)
		compareType := (intDateType-1)*24 + 12
		sql := "SELECT*FROM(SELECT count(*)AS a,company_id,nickname as company_name FROM x_attendance_record LEFT JOIN x_api_company ON x_api_company.id=x_attendance_record.company_id WHERE attendance_day>=date_format(date_sub(CURRENT_TIME(),interval+? HOUR),'%Y-%m-%d')AND onwork_time IS NOT NULL GROUP BY x_attendance_record.company_id)a LEFT JOIN(SELECT count(*)AS b,x_attendance_record.company_id FROM x_attendance_record LEFT JOIN x_api_role ON x_api_role.id=x_attendance_record.role_id WHERE attendance_day>=date_format(date_sub(CURRENT_TIME(),interval+? HOUR),'%Y-%m-%d')AND onwork_time IS NOT NULL AND role_code='seller_leader' GROUP BY x_attendance_record.company_id)b ON a.company_id=b.company_id LEFT JOIN(SELECT count(*)AS c,x_attendance_record.company_id FROM x_attendance_record LEFT JOIN x_api_role ON x_api_role.id=x_attendance_record.role_id WHERE onwork_time>=date_format(date_sub(CURRENT_TIME(),interval+? HOUR),'%Y-%m-%d')AND onwork_time IS NOT NULL AND role_code='seller' GROUP BY x_attendance_record.company_id)c ON c.company_id=a.company_id order by a desc limit ?"

		db.Raw(sql, compareType, compareType, compareType, limit).Scan(&attInfo)

		for _, value := range attInfo {
			chart = append(chart, Chart{value.CompanyName, value.A, "考勤人数"})
			chart = append(chart, Chart{value.CompanyName, value.B, "业务经理"})
			chart = append(chart, Chart{value.CompanyName, value.C, "酒水促销员"})
		}
		data, _ := json.Marshal(chart)
		rd.Set(redisKey, data, redisCacheTime*time.Second)

	}

	c.JSON(200, chart)
}

func SumPrice(c *gin.Context) {
	//var logger = util.Log()
	dateType := c.Param("type")
	limit, _ := strconv.Atoi(c.Param("limit"))
	var redisKey = "cache9-" + dateType + "-" + c.Param("limit")
	var db = model.DB
	var rd = cache.RedisClient
	var result []PriceSum
	var ext = rd.Exists(redisKey)
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		//logger.Info(jsonStr.Val())
		_ = json.Unmarshal([]byte(jsonStr.Val()), &result)

	} else {

		db.Raw("SELECT a.*,b.nickname FROM(SELECT SUM(price)/100 AS sum_price,company_id FROM x_user_ticket WHERE used_time>=DATE_SUB(NOW(),INTERVAL ? DAY)AND status=1 AND source_type=1 AND price>0 GROUP BY company_id)a LEFT JOIN x_api_company b ON a.company_id=b.id ORDER BY sum_price DESC LIMIT ?", dateType, limit).Scan(&result)

		data, _ := json.Marshal(result)
		rd.Set(redisKey, data, redisCacheTime*time.Second)

	}

	c.JSON(200, result)
}

func CompanyRoom(c *gin.Context) {
	//var logger = util.Log()
	dateType := c.Param("type")
	startRoomCount, _ := strconv.Atoi(c.Param("startRoomCount"))
	endRoomCount, _ := strconv.Atoi(c.Param("endRoomCount"))
	limit, _ := strconv.Atoi(c.Param("limit"))
	var redisKey = "cache10-" + dateType + "-" + c.Param("limit") + "-" + c.Param("startRoomCount") + "-" + c.Param("endRoomCount")
	var db = model.DB
	var rd = cache.RedisClient
	var result []CompanyRoomInfo
	var ext = rd.Exists(redisKey)
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		//logger.Info(jsonStr.Val())
		_ = json.Unmarshal([]byte(jsonStr.Val()), &result)

	} else {

		db.Raw("SELECT s.id,s.roomCount as room_count,s.n,s.avg_room,s.avg_up_user,s.avg_up_user/s.avg_room AS rjv FROM(SELECT company.id,e.roomCount,company.nickname AS n,CASE WHEN d.rooms IS NULL THEN 0 ELSE d.rooms END AS avg_room,CASE WHEN c.upUser IS NULL THEN 0 ELSE c.upUser END AS avg_up_user FROM x_api_company company LEFT JOIN(SELECT count(1)AS upUser,company_id FROM x_company_room_up_stage WHERE STATUS=0 AND add_time>DATE_SUB(NOW(),INTERVAL ? DAY)GROUP BY company_id)c ON c.company_id=company.id LEFT JOIN(SELECT count(1)AS rooms,company_id FROM x_company_room_consume_history WHERE add_time>DATE_SUB(NOW(),INTERVAL ? DAY)GROUP BY company_id)d ON d.company_id=company.id LEFT JOIN(SELECT count(1)AS roomCount,company_id FROM x_company_room WHERE status=0 AND is_delete=0 GROUP BY company_id)e ON e.company_id=company.id WHERE company.STATUS=0)s WHERE  s.roomCount > ? and s.roomCount< ? ORDER BY avg_room DESC LIMIT ?", dateType, dateType, startRoomCount, endRoomCount, limit).Scan(&result)

		data, _ := json.Marshal(result)
		rd.Set(redisKey, data, redisCacheTime*time.Second)

	}

	c.JSON(200, result)
}

func CompanyUpUser(c *gin.Context) {
	//var logger = util.Log()
	dateType := c.Param("type")
	limit, _ := strconv.Atoi(c.Param("limit"))
	var redisKey = "cache11-" + dateType + "-" + c.Param("limit")
	var db = model.DB
	var rd = cache.RedisClient
	var result []CompanyUpUserInfo
	var ext = rd.Exists(redisKey)
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		//logger.Info(jsonStr.Val())
		_ = json.Unmarshal([]byte(jsonStr.Val()), &result)

	} else {

		db.Raw("SELECT company.nickname,c.upUser AS up_user_count FROM x_api_company company LEFT JOIN(SELECT count(1)AS companyUser,company_id FROM x_api_user_company WHERE disabled NOT IN(2,4)AND STATUS IN(0,1)GROUP BY company_id)a ON a.company_id=company.id LEFT JOIN(SELECT count(1)AS attUser,company_id FROM x_attendance_record WHERE onwork_time IS NOT NULL AND onwork_time>DATE_SUB(NOW(),INTERVAL ? DAY)GROUP BY company_id)b ON b.company_id=company.id LEFT JOIN(SELECT count(1)AS upUser,company_id FROM x_company_room_up_stage WHERE STATUS=0 AND add_time>DATE_SUB(NOW(),INTERVAL ? DAY)GROUP BY company_id)c ON c.company_id=company.id WHERE company.STATUS=0 ORDER BY upuser DESC LIMIT ?", dateType, dateType, limit).Scan(&result)

		data, _ := json.Marshal(result)
		rd.Set(redisKey, data, redisCacheTime*time.Second)

	}

	c.JSON(200, result)
}

func DuplicateCompanyUser(c *gin.Context) {
	//var logger = util.Log()
	dateType := c.Param("type")
	var redisKey = "cache12-" + dateType
	var db = model.DB
	var rd = cache.RedisClient
	var result []DuplicateCompanyUserInfo
	var ext = rd.Exists(redisKey)
	if ext.Val() == 1 {
		var jsonStr = rd.Get(redisKey)
		//logger.Info(jsonStr.Val())
		_ = json.Unmarshal([]byte(jsonStr.Val()), &result)

	} else {

		db.Raw("SELECT v.company_name AS\"area\",a.y AS\"pv\",NOW()AS\"attribute\" FROM x_api_company v,(SELECT COMPANY_ID,count(*)AS Y FROM x_api_user_company WHERE creat_date>date_sub(curdate(),INTERVAL ? DAY)GROUP BY USER_ID)a WHERE v.id=a.COMPANY_ID AND a.company_id>12 AND a.Y>1 ORDER BY a.y DESC", dateType).Scan(&result)

		data, _ := json.Marshal(result)
		rd.Set(redisKey, data, redisCacheTime*time.Second)

	}

	c.JSON(200, result)
}

func SqlGet(sql string, dateType string, sum chan float64) {
	var db = model.DB
	var sum1 float64
	db.Raw(sql, dateType).Scan(&sum1)
	sum <- sum1
	//close(sum)
}
func SqlAttGet(sql string, dateType int, sum chan []AttInfo) {
	var db = model.DB
	var attList []AttInfo
	db.Raw(sql, dateType).Scan(&attList)
	sum <- attList
}
func round(x float64) int {
	return int(math.Floor(x + 0/5))
}

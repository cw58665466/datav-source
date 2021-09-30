package service

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"singo/model"
	"singo/serializer"
	"time"
)

// UserLoginService 管理用户登录的服务
type JwtLoginService struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// Login 用户登录函数
func (service *JwtLoginService) JwtLogin(c *gin.Context) serializer.Response {
	var user model.XPoliceUser
	fmt.Println(service)

	if err := model.DB.Where("email = ?", service.Email).First(&user).Error; err != nil {
		return serializer.ParamErr("账号或密码错误", nil)
	}

	if user.CheckNewPassword(service.Password) == false {
		return serializer.ParamErr("账号或密码错误", nil)
	}

	// 生成Token
	tokenString, _ := GenToken(user.ID)

	return serializer.LoginSuccess(tokenString, nil)
}

// MyClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录一个username字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type MyClaims struct {
	UserId uint `json:"userId"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("夏天夏天悄悄过去")

// GenToken 生成JWT
func GenToken(userId uint) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		userId, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "my-project",                               // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func Logout(tokenString string) {

}

func CurrentUser(c *gin.Context) model.XPoliceUser {

	userId := c.MustGet("userId")
	if userId != nil {
		user, _ := model.PoliceUser(userId)
		return user
	}
	return model.XPoliceUser{}
}

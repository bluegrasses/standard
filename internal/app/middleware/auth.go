package middleware

import (
	"net/http"
	"standard/internal/app/dao"
	"standard/internal/app/model"
	"standard/internal/global"
	"standard/pkg/tools"
	"time"

	"github.com/gin-gonic/gin"

	jwt "github.com/appleboy/gin-jwt/v2"
)

var identityKey = global.Conf.System.IdentityKey

//定义auth中间件
func InitAuth() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "standard",
		Key:         []byte("secret key"),
		Timeout:     time.Hour * 5,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		//登录②
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(map[string]interface{}); ok {
				return jwt.MapClaims{
					"user": v["user"],
				}
			}
			return jwt.MapClaims{}
		},
		//token 验证①
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return map[string]interface{}{
				"user": claims["user"],
			}
		},
		//登录①
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login model.Login
			if err := c.ShouldBind(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			u, err := dao.CheckAuth(login.Username, login.Password)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			user := tools.StructToJson(model.Claim{
				Id:       u.Id,
				Username: u.Username,
				Avatar:   u.Avatar,
				RoleId:   u.RoleId,
				RoleName: u.Role.Name,
			})
			return map[string]interface{}{
				"user": user,
			}, nil
		},
		//token 验证②
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(map[string]interface{}); ok {
				var user model.Claim
				// 将用户json转为结构体
				tools.JsonI2Struct(v["user"], &user)
				//这里推荐存入struct,方便获取值
				c.Set("user", &user)
				return true
			}
			return false
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			global.Logger.Debug(message)
			c.JSON(http.StatusOK, gin.H{
				"code": 405,
				"Msg":  "jwt 校验失败",
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"xtuOj/helper"
)

func CheckUser() gin.HandlerFunc {

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		userClaim, err := helper.AnalyseToken(auth)
		// 如果解析错误
		if err != nil {
			c.JSON(http.StatusOK, &gin.H{
				"code": -1,
				"msg":  "用户错误",
			})
			log.Println("middleware checkAuthAdmin 解析用户token错误: ", err)
			c.Abort()
			return
		}

		c.Set("user", userClaim)
		c.Next()
	}
}

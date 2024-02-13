package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
	"xtuOj/define"
	"xtuOj/helper"
	"xtuOj/models"
)

// GetUserDetail
// @Tags 公共方法
// @Summary 获取用户信息
// @Param identity query string false "identity"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context) {
	identity := c.Query("identity")

	if identity == "" {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "请输入用户id",
		})
		return
	}
	data, err := models.GetUserDetail(identity)

	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "数据库查询错误",
		})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"data": data,
	})

}

// Login
// @Tags 公共方法
// @Summary 用户登录
// @Param username formData string false "用户名"
// @Param password formData string false "密码"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /login [post]
func Login(c *gin.Context) {
	// 获取参数
	password := c.PostForm("password")
	userName := c.PostForm("username")

	if password == "" || userName == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "need username or password",
		})
		return
	}

	// 将密码转化成md5
	password = helper.GetMd5(password)
	data, e := models.Login(userName, password)

	// 进行md5的错误处理
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, &gin.H{
				"code": -1,
				"msg":  "不存在此用户",
			})
			return
		}
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})
		return
	}

	// 验证成功转化成token
	tokenString, e := helper.GetToken(data.Name, data.Identity, data.IsAdmin)

	if e != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "token验证错误",
		})
		return
	}

	// 返回tokenstring
	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"data": gin.H{
			"token": tokenString,
		},
	})
}

// SendMailCode
// @Tags 公共方法
// @Summary 发送邮箱验证码
// @Param email formData string true "电子邮箱"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /send-code [post]
func SendMailCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})

		return
	}

	//TODO:生成code
	code := helper.GetRandCode()

	// 设置五分钟之后过期
	models.RDB.Set(email, code, time.Second*300)

	err := helper.SendEmailCode(email, code)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "send email error",
		})
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"msg":  "send email successful",
	})
}

// Register
// @Tags 公共方法
// @Summary 用户注册
// @Param mail formData string true "电子邮箱"
// @Param name formData string true "用户名"
// @Param password formData string true "密码"
// @Param code formData string true "验证码"
// @Param phone formData string false "手机"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /register [post]
func Register(c *gin.Context) {
	// 获取相关参数
	mail := c.PostForm("mail")
	name := c.PostForm("name")
	code := c.PostForm("code")
	password := c.PostForm("password")
	password = helper.GetMd5(password)
	phone := c.PostForm("phone")

	if mail == "" || name == "" || code == "" || password == "" {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}

	// 验证code是否正确
	trueCode, err := models.RDB.Get(mail).Result()
	if err != nil && trueCode != code {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "验证码错误",
		})
		return
	}

	// 判断用户是否存在
	ok, err := models.IsUserExist(mail)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "系统服务器错误",
		})
		log.Println("check user is existed error!", err)
		return
	}
	if !ok {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "用户已经存在",
		})
		return
	}

	// 插入数据库
	identity, err := models.Register(name, password, phone, mail)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "系统服务器错误",
		})
		log.Println("insert database error!", err)
		return
	}

	// 获取用户token
	usertoken, err := helper.GetToken(name, identity, 0)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "无法验证用户身份",
		})
		log.Println("userToken error:", err)
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"data": gin.H{
			"token": usertoken,
		},
	})
}

// GetRankList
// @Tags 公共方法
// @Summary 用户排行榜
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /rank-list [get]
func GetRankList(c *gin.Context) {
	// 获取参数
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("page need a int param :", err)
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("siz need a int param:", err)
	}
	offset := (page - 1) * size

	// 查询数据库
	data, count, err := models.GetRankList(offset, size)

	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})

		log.Println("查询数据失败:", err)
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"data": gin.H{
			"list":  data,
			"count": count,
		},
	})
}

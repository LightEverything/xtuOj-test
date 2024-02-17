package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"strconv"
	"xtuOj/define"
	"xtuOj/execode"
	"xtuOj/helper"
	"xtuOj/models"
)

// GetSubmitList
// @Tags 公共方法
// @Summary 提交列表
// @Param page query int false "请输入当前页数,默认第一页"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity "
// @Param user_identity query string false "user_identity"
// @Param status query int false "status"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	// 获取相关参数
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("page need a int param :", err)
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("siz need a int param:", err)
	}

	status, err := strconv.Atoi(c.DefaultQuery("status", define.DefaultStatus))
	if err != nil {
		log.Println("status need a int param", err)
	}

	offset := (page - 1) * size
	// 问题identity
	problemIdentity := c.Query("problem_identity")
	// user的identity
	userIdentity := c.Query("user_identity")

	// 调用model接口
	data, count, err := models.GetSubmitList(offset, size, problemIdentity, userIdentity, status)

	if err != nil {
		log.Println("submit-detail db error :", err)
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "server error",
		})
		return
	}
	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"data": gin.H{
			"count": count,
			"list":  data,
		},
	})
}

// Submit
// @Tags 用户私有方法
// @Summary 代码提交
// @Param Authorization header string true "token"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /user/submit [post]
func Submit(c *gin.Context) {
	// 获取参数
	problem_identity := c.Query("problem_identity")
	code, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "获取body失败",
		})

		log.Println("get submit body error :", err)
		return
	}

	// 创建代码保存路径
	path, err := helper.SaveCode(code)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})

		log.Println("saveCode error : ", err)
		return
	}

	// 获取上下文信息
	identity, _ := helper.GetUuid()
	user, _ := c.Get("user")
	uc := user.(*helper.UserClaims)

	submitSt := models.Submit{
		Identity:        identity,
		ProblemIdentity: problem_identity,
		UserIdentity:    uc.Identity,
		Path:            path,
	}

	// 执行代码判断
	data, err := models.GetProblemDetail(problem_identity)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})
		log.Println("getTestCase error:", err)
		return
	}

	er := execode.ExeCode(context.Background(), data, path)

	submitSt.Status = er.Status

	// 创建新的提交记录
	if err := models.CreateSubmission(&submitSt, data, problem_identity); err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})
		log.Println("create submission error:", err)
		return
	}

	// 返回对应信息
	c.JSON(http.StatusOK, &gin.H{
		"code":   200,
		"msg":    er.Msg,
		"status": er.Status,
	})
}

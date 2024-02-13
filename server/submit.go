package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"xtuOj/define"
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

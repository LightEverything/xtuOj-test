package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"xtuOj/define"
	"xtuOj/models"
)

// GetCategoryList
// @Tags 私有方法
// @Summary 分类列表
// @Param Authorization header string true "token"
// @Param page query int false "请输入当前页数,默认第一页"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /category-list [get]
func GetCategoryList(c *gin.Context) {
	// 获取相关参数
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("page need a int param :", err)
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("siz need a int param:", err)
	}

	keyword := c.Query("keyword")

	// 设置偏移量
	offset := (page - 1) * size

	data, count, e := models.GetCategoryList(offset, size, keyword)
	if e != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})
		log.Println("getcategorylist 数据服务错误：", err)
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

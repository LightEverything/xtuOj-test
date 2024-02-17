package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"xtuOj/define"
	"xtuOj/helper"
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
// @Router /admin/category-list [get]
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

// CreateCategory
// @Tags 私有方法
// @Summary 创建分类列表
// @Param Authorization header string true "token"
// @Param name formData string true "分类名称"
// @Param parentId formData int false "parentId"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /admin/category-create [post]
func CreateCategory(c *gin.Context) {
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("parentId"))

	identity, err := helper.GetUuid()
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "系统服务错误",
		})

		log.Println("CreateCategory error :", err)
		return
	}

	// 根据参数构建data
	data := models.Category{
		Identity: identity,
		Name:     name,
		Parentid: parentId,
	}

	// 在数据库插入相关数据
	err = models.CreateCategory(data)

	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "系统服务错误",
		})

		log.Println("CreateCategory DB error :", err)
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"msg":  "创建成功",
	})
}

// DeleteCategory
// @Tags 私有方法
// @Summary 修改分类列表
// @Param Authorization header string true "token"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /admin/category-delete [delete]
func DeleteCategory(c *gin.Context) {
	identity := c.Query("identity")
	count, err := models.GetProblemsOfCategoryCount(identity)

	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})
		log.Println("deletecategory 数据库错误：", err)
		return
	}

	// 如果该分类下面关联相关的问题，则无法删除
	if count > 0 {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "存在相关分类问题, 无法删除",
		})
		return
	}

	err = models.DeleteCategory(identity)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})
		log.Println("deletecategory 数据库错误：", err)
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// UpdateCategory
// @Tags 私有方法
// @Summary 修改分类列表
// @Param Authorization header string true "token"
// @Param parentId formData int false "parentId"
// @Param identity formData string true "identity"
// @Param name formData string true "分类名称"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /admin/category-update [put]
func UpdateCategory(c *gin.Context) {

	// 获取参数
	identity := c.PostForm("identity")
	name := c.PostForm("name")
	parentId := c.PostForm("parentId")

	var parentIdint int

	if parentId != "" {
		tmp, err := strconv.Atoi(parentId)
		if err != nil {
			c.JSON(http.StatusOK, &gin.H{
				"code": -1,
				"msg":  "参数错误",
			})
			return
		}

		parentIdint = tmp
	}

	data := models.Category{
		Identity: identity,
		Name:     name,
		Parentid: parentIdint,
	}

	// 链接数据库
	err := models.UpdateCategory(data)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})
		log.Println("updateCategory 更新数据库错误：", err)
		return
	}

	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"msg":  "修改成功",
	})
}

package server

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"xtuOj/define"
	"xtuOj/helper"
	"xtuOj/models"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "请输入当前页数,默认第一页"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /problem-list [get]
func GetProblemList(c *gin.Context) {
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
	categoryId := c.Query("category_id")

	// 设置偏移量
	offset := (page - 1) * size

	data, count, err := models.GetProblemList(offset, size, keyword, categoryId)

	if err != nil {
		log.Println("DB get data failure:", err)
	}

	// 返回json
	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"data": gin.H{
			"list":  data,
			"count": count,
		},
	})

}

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "identity"
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")

	if identity == "" {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "没有填写信息",
		})
		return
	}

	data, err := models.GetProblemDetail(identity)

	// 错误处理
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, &gin.H{
				"code": -1,
				"msg":  "没有找到对应记录",
			})
			return
		}
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务端错误",
		})
	} else {
		// 若未发生错误，则返回相关信息
		c.JSON(http.StatusOK, &gin.H{
			"code": 200,
			"data": data,
		})
	}

}

// 处理categoryids
func dealCategoryIds(categoryIds []string, data *models.Problem) ([]*models.ProblemCategory, error) {

	categoryArray := make([]*models.ProblemCategory, 0)

	for _, v := range categoryIds {
		cid, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		// 扩展一个新的problem_category
		categoryArray = append(categoryArray, &models.ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(cid),
		})
	}

	return categoryArray, nil
}

// 处理testcaseid
func dealTestCases(testCases []string, data *models.Problem) ([]*models.TestCase, error) {
	testCaseArray := make([]*models.TestCase, 0)

	for _, v := range testCases {
		dataMap := make(map[string]string)
		// 解析testcases的json值
		if err := json.Unmarshal([]byte(v), &dataMap); err != nil {
			return nil, err
		}

		// 判读输入输出是否存在
		inputData, ok := dataMap["input"]
		if !ok {
			return nil, errors.New("无输入参数")
		}
		outputData, ok := dataMap["output"]
		if !ok {
			return nil, errors.New("无输出参数")
		}

		// 获取输入输出唯一标识符号
		identity, err := helper.GetUuid()
		if err != nil {
			return nil, err
		}

		testCaseArray = append(testCaseArray, &models.TestCase{
			Identity:        identity,
			ProblemIdentity: data.Identity,
			Output:          outputData,
			Input:           inputData,
		})

	}

	return testCaseArray, nil
}

// CreateProblem
// @Tags 私有方法
// @Summary 创建问题
// @Param Authorization header string true "token"
// @Param title formData string true "标题"
// @Param content formData string true "内容"
// @Param max_runtime formData int true "max_runtime"
// @Param max_mem formData int true "max_mem"
// @Param category_ids formData []string false "标签" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /admin/problem-create [post]
func CreateProblem(c *gin.Context) {
	// 获取相关参数
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))

	testCases := c.PostFormArray("test_cases")
	categoryIds := c.PostFormArray("category_ids")

	if title == "" || content == "" || len(categoryIds) == 0 || len(testCases) == 0 {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "参数设置错误",
		})
		log.Println("参数设置失败(有必填参数为空)")
		return
	}

	// 获取uuid
	identity, err := helper.GetUuid()

	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})
		log.Println("uuid生成失败")
	}

	data := models.Problem{
		Model:      gorm.Model{},
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMem:     maxMem,
	}

	// 处理 testCases 与 categoryIds
	categoryArray, err := dealCategoryIds(categoryIds, &data)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		log.Println("获取categoryArray错误：", err)
		return
	}
	data.ProblemCategorys = categoryArray

	testCaseArray, err := dealTestCases(testCases, &data)

	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		log.Println("获取testCaseArray错误：", err)
		return
	}
	data.TestCases = testCaseArray
	// 创建一个新problem
	err = models.CreateProblem(&data)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "系统服务器错误",
		})

		log.Println("CreateProblem数据库链接错误:", err)
		return
	}

	// 若成功创建，则返回identity
	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"data": gin.H{
			"identity": data.Identity,
		},
	})
}

// UpdateProblem
// @Tags 私有方法
// @Summary 修改问题
// @Param Authorization header string true "token"
// @Param title formData string true "标题"
// @Param identity formData string true "identity"
// @Param content formData string true "内容"
// @Param max_runtime formData int true "max_runtime"
// @Param max_mem formData int true "max_mem"
// @Param category_ids formData []string false "标签" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":"200", "data":""}"
// @Router /admin/problem-update [put]
func UpdateProblem(c *gin.Context) {
	// 获取相关参数
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	identity := c.PostForm("identity")

	testCases := c.PostFormArray("test_cases")
	categoryIds := c.PostFormArray("category_ids")

	if title == "" || content == "" || len(categoryIds) == 0 || len(testCases) == 0 {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "参数设置错误",
		})
		log.Println("参数设置失败(有必填参数为空)")
		return
	}

	// 根据参数创建一个对应的Problem struct
	// 1. 获取id
	dataid, err := models.GetProblemId(identity)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})

		log.Println("UpdateProblem database error:", err)
		return
	}

	// 2. 填充基本参数
	data := models.Problem{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMem:     maxMem,
	}

	data.ID = dataid

	// 3. 填充数组参数
	categoryArray, err := dealCategoryIds(categoryIds, &data)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		log.Println("获取categoryArray错误：", err)
		return
	}
	data.ProblemCategorys = categoryArray

	testCaseArray, err := dealTestCases(testCases, &data)

	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "参数错误",
		})
		log.Println("获取testCaseArray错误：", err)
		return
	}
	data.TestCases = testCaseArray

	// 链接数据库更新问题
	err = models.UpdateProblem(&data)
	if err != nil {
		c.JSON(http.StatusOK, &gin.H{
			"code": -1,
			"msg":  "服务器错误",
		})
		log.Println("updateProblem database error:", err)
		return
	}

	//返回正确值
	c.JSON(http.StatusOK, &gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

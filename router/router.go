package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "xtuOj/docs"
	"xtuOj/middleware"
	"xtuOj/server"
)

func Router() *gin.Engine {
	r := gin.Default()

	/*
			公有方法
		******************************************************
	*/

	// 测试路由
	r.GET("/ping", server.Ping)

	// problem路由
	r.GET("/problem-list", server.GetProblemList)
	r.GET("/problem-detail", server.GetProblemDetail)

	// submit路由
	r.GET("/submit-list", server.GetSubmitList)

	//user路由
	r.GET("/user-detail", server.GetUserDetail)
	r.POST("/login", server.Login)
	r.POST("/register", server.Register)
	r.POST("/send-code", server.SendMailCode)
	r.GET("/rank-list", server.GetRankList)

	// swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	/*
		私有方法
		*********************************************************
	*/

	GroupAdmin := r.Group("/admin", middleware.CheckAuthAdmin())
	{
		// 创建问题
		GroupAdmin.POST("/problem-create", server.CreateProblem)
		GroupAdmin.PUT("/problem-update", server.UpdateProblem)
		// 获取相关列表
		GroupAdmin.GET("/category-list", server.GetCategoryList)
		// 删除相关分类
		GroupAdmin.DELETE("/category-delete", server.DeleteCategory)
		// 修改相关分类
		GroupAdmin.PUT("/category-update", server.UpdateCategory)
		// 创建相关分类
		GroupAdmin.POST("/category-create", server.CreateCategory)
	}

	GroupUser := r.Group("/user", middleware.CheckUser())
	{
		GroupUser.POST("/submit", server.Submit)
	}

	return r
}

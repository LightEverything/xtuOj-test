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

	r.POST("/problem-create", middleware.CheckAuthAdmin(), server.CreateProblem)
	r.GET("/category-list", middleware.CheckAuthAdmin(), server.GetCategoryList)

	return r
}

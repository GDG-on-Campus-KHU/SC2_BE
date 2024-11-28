package routes

import (
	"github.com/GDG-on-Campus-KHU/SC2_BE/controllers"
	"github.com/GDG-on-Campus-KHU/SC2_BE/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/GDG-on-Campus-KHU/SC2_BE/db"
)

func Routes() *gin.Engine {
	router := gin.Default()

	// 스웨거 라우트
	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// map routes group
	mapRoutes := router.Group("/api/map")
	{
		mapRoutes.GET("/search", controllers.NaverSearchHandler)
		mapRoutes.GET("/navigation", controllers.GetNavigateHandler)
		mapRoutes.GET("/delete", db.DeleteAllDocument)
	}

	// 재난 알림 문자 라우트
	router.GET("/api/disaster-messages", controllers.GetDisasterMessagesHandler)

	// AI 모델에 재난 문자 전송하고 응답 받는 라우트
	router.POST("/api/ai/disaster-messages", controllers.SendDisasterMessageController)

	// 알림 라우트
	router.POST("/api/send", controllers.SendNotificationHandler)

	return router
}

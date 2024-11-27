package routes

import(
	"github.com/gin-gonic/gin"
	"github.com/GDG-on-Campus-KHU/SC2_BE/controllers"
	"github.com/GDG-on-Campus-KHU/SC2_BE/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes() *gin.Engine{
	router := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// map routes group
	mapRoutes := router.Group("/api/map")
	{
		mapRoutes.GET("/search", controllers.QuerySearch)
		mapRoutes.GET("/navigation", controllers.GetNavigateHandler)
	}
	return router
}
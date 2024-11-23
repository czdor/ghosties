package routes

import (
	"core/internal/middleware"
	"core/internal/server/controllers"

	"core/internal/ws"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// REST API
	apiv1 := r.Group("/api/v1")
	{
		userGroup := apiv1.Group("/user")
		{
			userGroup.POST("/signup", controllers.Signup)
			userGroup.POST("/login", controllers.Login)
			userGroup.GET("/refresh", middleware.Authenticate, controllers.Refresh)
			userGroup.GET("/protected", middleware.Authenticate, controllers.Protected) // ! this is a mock ep
		}

		roomGroup := apiv1.Group("/rooms")
		{
			roomGroup.GET("", controllers.GetRooms)
		}

	}

	// Docs WebSocket API
	r.GET("/ws", ws.HandleWebSocket)
	r.POST("/ws", ws.HandleWebSocket)

	// WebSocket API Docs
	r.Static("/wsdocs", "./docs")

	// Docs REST API

	// docs.SwaggerInfo.Title = "my api title"
	// docs.SwaggerInfo.Description = "my description"
	// docs.SwaggerInfo.Version = "1.0"

	// r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
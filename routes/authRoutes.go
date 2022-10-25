package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/rwiteshbera/orbit/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", controllers.Signup)
	// incomingRoutes.POST("/login", controllers.Login)
}

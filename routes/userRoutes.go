package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/rwiteshbera/orbit/controllers"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// incomingRoutes.Use(middlewares.Authenticate()) // Check whether the routes are authenticated or not
	incomingRoutes.GET("/user", controllers.GetUser())
}

package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/rwiteshbera/orbit/controllers"
	middlewares "github.com/rwiteshbera/orbit/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middlewares.Authenticate())                   // Check whether the routes are authenticated or not
	incomingRoutes.GET("/user", controllers.GetUser())               //  /user or /user?name=user123
	incomingRoutes.POST("/user/edit", controllers.EditUser())        // Edit account details
	incomingRoutes.POST("/user/delete", controllers.DeleteAccount()) // Delete user account
}

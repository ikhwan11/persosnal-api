package routes

import (
	"my-personal-web/api/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Users endpoint
	router.POST("/users", controller.CreateUserHandler)
}

package routes

import (
	"my-personal-web/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	// Users endpoint
	router.POST("/users", controller.CreateUserHandler)
	router.GET("/users", controller.GetUsersHandler)
	router.GET("/users/:slug", controller.GetUserBySlugHandler)
	router.PUT("/users/:slug", controller.UpdateUserHandler)
}

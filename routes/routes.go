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

	// cv endpoint
	router.GET("/curriculum-vitae", controller.GetCVHandler)
	router.POST("/curriculum-vitae", controller.CreateCVHandler)
	router.PUT("/curriculum-vitae/:slug", controller.UpdateCVHandler)
}

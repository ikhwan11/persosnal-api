package routes

import (
	"my-personal-web/controller"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Users endpoint
	r.POST("/users", controller.CreateUserHandler)
	r.GET("/users", controller.GetUsersHandler)
	r.GET("/users/:slug", controller.GetUserBySlugHandler)
	r.PUT("/users/:slug", controller.UpdateUserHandler)

	// cv endpoint
	r.GET("/curriculum-vitae", controller.GetCVHandler)
	r.POST("/curriculum-vitae", controller.CreateCVHandler)
	r.PUT("/curriculum-vitae/:slug", controller.UpdateCVHandler)
	// skills endpoint
	r.POST("/skills/:cv_id", controller.CreateSkillsHandler)
	r.PUT("/skills/:skill_id", controller.UpdateSkillsHandler)
	r.GET("/skills/:cv_id/hard", controller.GetHardSkillsHandler)
	r.GET("/skills/:cv_id/soft", controller.GetSoftSkillsHandler)
	r.GET("/skills/:cv_id/tools", controller.GetToolsSkillsHandler)
	r.DELETE("/skills/:id", controller.DeleteSkillHandler)
	// educations endpoint
	r.POST("/educations/:cv_id", controller.CreateEducationsHandler)
	r.PUT("/educations/:edu_id", controller.UpdateEducationsHandler)
	r.GET("/educations/:cv_id", controller.GetEducationsHandler)
	r.DELETE("/educations/:id", controller.DeleteEducationHandler)
}

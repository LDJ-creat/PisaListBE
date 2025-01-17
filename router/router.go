package router

import (
	v1 "github.com/PisaListBE/api/v1"
	"github.com/PisaListBE/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	// 使用 CORS 中间件
	r.Use(middleware.Cors())

	api := r.Group("/api/v1")
	{
		// 公开路由 - 不需要 JWT 验证
		api.POST("/register", v1.Register)
		api.POST("/login", v1.Login)
		api.GET("/wishes/community", v1.GetCommunityWishes)
		api.GET("/wishes/random", v1.GetRandomWish)

		// 需要验证的路由组
		auth := api.Group("")
		auth.Use(middleware.JWT())
		{
			// 任务相关路由
			auth.POST("/tasks", v1.CreateTask)
			auth.GET("/tasks/today", v1.GetTodayTasks)
			auth.GET("/tasks/timeline", v1.GetTaskTimeline)
			auth.PUT("/tasks/:id", v1.UpdateTask)
			auth.DELETE("/tasks/:id", v1.DeleteTask)
			auth.PUT("/tasks/:id/complete", v1.CompleteTask)
			auth.PUT("/tasks/importance", v1.UpdateTasksImportance)

			// 需要验证的心愿相关路由
			auth.POST("/wishes", v1.CreateWish)
			auth.GET("/wishes", v1.GetUserWishes)
			auth.PUT("/wishes/:id", v1.UpdateWish)
			auth.DELETE("/wishes/:id", v1.DeleteWish)
			auth.POST("/wishes/:id/share", v1.ShareWish)
		}
	}
}

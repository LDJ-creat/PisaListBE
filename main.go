package main

// @title PisaList API
// @version 1.0
// @description PisaList 是一个待办事项和心愿清单管理系统
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

import (
	"log"

	v1 "github.com/PisaListBE/api/v1"
	_ "github.com/PisaListBE/docs"
	"github.com/PisaListBE/internal/middleware"
	"github.com/PisaListBE/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}
}

func main() {
	gin.SetMode(viper.GetString("server.mode"))
	r := gin.Default()

	r.Use(middleware.Cors())
	setupRoutes(r)

	port := viper.GetString("server.port")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server error: %s", err)
	}
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		// 公开路由
		api.POST("/register", v1.Register)
		api.POST("/login", v1.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.JWT())
		{
			// Task相关路由
			auth.POST("/tasks", v1.CreateTask)
			auth.DELETE("/tasks/:id", v1.DeleteTask)
			auth.PUT("/tasks/:id", v1.UpdateTask)
			auth.PUT("/tasks/:id/complete", v1.CompleteTask)
			auth.PUT("/tasks/:id/importance", v1.UpdateTaskImportance)
			auth.GET("/tasks/timeline", v1.GetTaskTimeline)
			auth.GET("/tasks/today", v1.GetTodayTasks)

			// Wish相关路由
			auth.POST("/wishes", v1.CreateWish)
			auth.DELETE("/wishes/:id", v1.DeleteWish)
			auth.PUT("/wishes/:id", v1.UpdateWish)
			auth.POST("/wishes/:id/share", v1.ShareWish)
			auth.GET("/wishes", v1.GetUserWishes)
			auth.GET("/wishes/community", v1.GetCommunityWishes)
			auth.GET("/wishes/random", v1.GetRandomWish)
		}
	}
}

package main

import (
	"github.com/PisaListBE/pkg/database"
	"github.com/PisaListBE/router"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	if err := database.InitGormDB(); err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	r := gin.Default()

	// 初始化路由
	router.InitRouter(r)

	// 启动服务器
	r.Run(":8080")
}

package database

import (
	"fmt"

	"github.com/PisaListBE/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func initSharedWishes() {
	var count int64
	GormDB.Model(&model.SharedWish{}).Count(&count)
	if count == 0 {
		wishes := []model.SharedWish{
			{Event: "珍惜时光，享受当下", Description: "春风若有怜花意，可否许我在少年"},
			{Event: "不忘初心，牢记使命", Description: "奋勇争先，不负韶华"},
			{Event: "Be confident all the time", Description: "仰天大笑出门去，我辈岂是蓬蒿人"},
			{Event: "来一场说走就走的旅行吧", Description: ""},
			{Event: "永远相信美好的事情即将发生", Description: "只要出发了，我们就在通往胜利的路上"},
			{Event: "己所不欲，勿施于人", Description: ""},
			{Event: "穷则兼善其身，达则兼济天下", Description: "天下兴亡，匹夫有责"},
			{Event: "永远积极向上，永远热泪盈眶，永远豪情满怀，永远坦坦荡荡", Description: ""},
			{Event: "把我的技能包点满☺☺☺", Description: "不忘初心，牢记使命"},
			{Event: "柴米油盐皆是诗，无灾无难是最重", Description: "如果快乐太难，祝你我都平平安安"},
		}

		for _, wish := range wishes {
			GormDB.Create(&wish)
		}
	}
}

func InitGormDB() error {
	dsn := "root:268968&&ABc@tcp(localhost:3306)/pisa_list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 先创建数据库（如果不存在）
	err = db.Exec("CREATE DATABASE IF NOT EXISTS pisa_list").Error
	if err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&model.Task{}, &model.Wish{}, &model.SharedWish{}, &model.User{})
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	GormDB = db

	// 初始化心愿社区数据
	initSharedWishes()

	return nil
}

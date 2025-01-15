package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PisaListBE/internal/model"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// 预设的心愿列表
var defaultWishes = []model.SharedWish{
	{
		Event:          "环游世界",
		Description:    "去不同的国家体验不同的文化",
		SharedByUserID: 1,
	},
	{
		Event:          "学习一门新语言",
		Description:    "掌握一门外语，开拓新的视野",
		SharedByUserID: 1,
	},
	{
		Event:          "参加马拉松",
		Description:    "通过长期训练完成一次马拉松比赛",
		SharedByUserID: 1,
	},
	{
		Event:          "学习烹饪",
		Description:    "掌握烹饪技能，为家人做美食",
		SharedByUserID: 1,
	},
	{
		Event:          "创办一家公司",
		Description:    "实现创业梦想，创造社会价值",
		SharedByUserID: 1,
	},
	{
		Event:          "写一本书",
		Description:    "记录自己的思考和经历",
		SharedByUserID: 1,
	},
	{
		Event:          "学习摄影",
		Description:    "用镜头记录生活中的美好瞬间",
		SharedByUserID: 1,
	},
	{
		Event:          "参与志愿服务",
		Description:    "为社会做出贡献，帮助他人",
		SharedByUserID: 1,
	},
	{
		Event:          "学习乐器",
		Description:    "培养音乐素养，丰富生活",
		SharedByUserID: 1,
	},
	{
		Event:          "种植一个花园",
		Description:    "亲手培育植物，感受生命的成长",
		SharedByUserID: 1,
	},
}

func InitDB() error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.dbname"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(viper.GetInt("database.max_idle_conns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("database.max_open_conns"))

	// 手动创建表结构
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
			created_at DATETIME(3) NULL,
			updated_at DATETIME(3) NULL,
			deleted_at DATETIME(3) NULL,
			user_id BIGINT UNSIGNED NOT NULL,
			event VARCHAR(256) NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			is_cycle BOOLEAN DEFAULT FALSE,
			description TEXT,
			importance_level INT DEFAULT 0,
			completed_date DATETIME NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`).Error
	if err != nil {
		return err
	}

	// 自动迁移其他表
	err = db.AutoMigrate(
		&model.User{},
		&model.Wish{},
		&model.SharedWish{},
	)
	if err != nil {
		return err
	}

	DB = db

	// 初始化心愿社区数据
	return initializeWishCommunity(db)
}

// 初始化心愿社区数据
func initializeWishCommunity(db *gorm.DB) error {
	// 检查是否已经有数据
	var count int64
	if err := db.Model(&model.SharedWish{}).Count(&count).Error; err != nil {
		return err
	}

	// 如果没有数据，则插入预设的心愿
	if count == 0 {
		for _, wish := range defaultWishes {
			if err := db.Create(&wish).Error; err != nil {
				return err
			}
		}
		log.Println("Successfully initialized wish community with default wishes")
	}

	return nil
}

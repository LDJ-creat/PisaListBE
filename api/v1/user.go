package v1

import (
	"fmt"
	"net/http"

	"github.com/PisaListBE/internal/model"
	"github.com/PisaListBE/pkg/database"
	"github.com/PisaListBE/pkg/jwt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @title PisaList User API
// @version 1.0
// @description 用户管理相关的API接口

// UserRequest 用户请求结构体
// @Description 用户注册请求的数据结构
type UserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32" example:"johndoe" description:"用户名，3-32个字符"`
	Password string `json:"password" binding:"required,min=6" example:"password123" description:"密码，最少6个字符"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com" description:"电子邮件地址"`
}

// @Summary 用户注册
// @Description 创建新用户账号
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserRequest true "用户注册信息"
// @Success 200 {object} object{token=string,user=object{id=integer,username=string,email=string}} "注册成功返回token和用户信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /register [post]
func Register(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查用户名是否已存在
	var existingUser model.User
	if err := database.GormDB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := database.GormDB.Create(&user).Error; err != nil {
		// 添加详细的错误日志
		fmt.Printf("创建用户失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建用户失败"})
		return
	}

	// 生成 token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// @Summary 用户登录
// @Description 用户登录并获取JWT令牌
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body object{username=string,password=string} true "登录凭证"
// @Success 200 {object} object{token=string,user=object{id=integer,username=string,email=string}} "登录成功返回token和用户信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "用户名或密码错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /login [post]
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required" example:"johndoe"`
		Password string `json:"password" binding:"required" example:"password123"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if err := database.GormDB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

package v1

import (
	"math/rand"
	"net/http"

	"github.com/PisaListBE/internal/model"
	"github.com/PisaListBE/pkg/database"
	"github.com/gin-gonic/gin"
)

// @title PisaList Wish API
// @version 1.0
// @description 心愿管理相关的API接口

// WishRequest 心愿请求结构体
// @Description 创建或更新心愿的请求数据结构
type WishRequest struct {
	Event       string `json:"event" binding:"required" example:"环游世界" description:"心愿内容"`
	Description string `json:"description" example:"想去看看世界的每个角落" description:"心愿详细描述"`
	IsCycle     bool   `json:"is_cycle" example:"false" description:"是否为循环心愿"`
}

// @Summary 创建心愿
// @Description 创建一个新的心愿
// @Tags wishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param wish body WishRequest true "心愿信息"
// @Success 200 {object} model.Wish "创建成功返回心愿信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /wishes [post]
func CreateWish(c *gin.Context) {
	userID := c.GetUint("userID")
	var req WishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	wish := model.Wish{
		UserID:      userID,
		Event:       req.Event,
		Description: req.Description,
		IsCycle:     req.IsCycle,
		IsShared:    false,
	}

	if err := database.GormDB.Create(&wish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建心愿失败"})
		return
	}

	c.JSON(http.StatusOK, wish)
}

// @Summary 删除心愿
// @Description 删除指定的心愿
// @Tags wishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "心愿ID"
// @Success 200 {object} map[string]string "删除成功"
// @Failure 404 {object} map[string]string "心愿不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /wishes/{id} [delete]
func DeleteWish(c *gin.Context) {
	userID := c.GetUint("userID")
	wishID := c.Param("id")

	var wish model.Wish
	if err := database.GormDB.Where("id = ? AND user_id = ?", wishID, userID).First(&wish).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "心愿不存在"})
		return
	}

	if err := database.GormDB.Delete(&wish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除心愿失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// @Summary 更新心愿
// @Description 更新指定心愿的信息
// @Tags wishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "心愿ID"
// @Param wish body WishRequest true "更新的心愿信息"
// @Success 200 {object} model.Wish "更新成功返回心愿信息"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 404 {object} map[string]string "心愿不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /wishes/{id} [put]
func UpdateWish(c *gin.Context) {
	userID := c.GetUint("userID")
	wishID := c.Param("id")

	var req WishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wish model.Wish
	if err := database.GormDB.Where("id = ? AND user_id = ?", wishID, userID).First(&wish).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "心愿不存在"})
		return
	}

	updates := map[string]interface{}{
		"event":       req.Event,
		"description": req.Description,
		"is_cycle":    req.IsCycle,
	}

	if err := database.GormDB.Model(&wish).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新心愿失败"})
		return
	}

	c.JSON(http.StatusOK, wish)
}

// @Summary 分享心愿
// @Description 将心愿分享到心愿社区
// @Tags wishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "心愿ID"
// @Success 200 {object} map[string]string "分享成功"
// @Failure 404 {object} map[string]string "心愿不存在"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /wishes/{id}/share [post]
func ShareWish(c *gin.Context) {
	userID := c.GetUint("userID")
	wishID := c.Param("id")

	var wish model.Wish
	if err := database.GormDB.Where("id = ? AND user_id = ?", wishID, userID).First(&wish).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "心愿不存在"})
		return
	}

	sharedWish := model.SharedWish{
		OriginalWishID: wish.ID,
		Event:          wish.Event,
		Description:    wish.Description,
		SharedByUserID: userID,
	}

	if err := database.GormDB.Create(&sharedWish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "分享心愿失败"})
		return
	}

	if err := database.GormDB.Model(&wish).Update("is_shared", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新心愿分享状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分享成功"})
}

// @Summary 获取用户心愿列表
// @Description 获取当前用户的所有心愿
// @Tags wishes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Wish "心愿列表"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /wishes [get]
func GetUserWishes(c *gin.Context) {
	userID := c.GetUint("userID")

	var wishes []model.Wish
	if err := database.GormDB.Where("user_id = ?", userID).Find(&wishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取心愿列表失败"})
		return
	}

	c.JSON(http.StatusOK, wishes)
}

// @Summary 获取心愿社区列表
// @Description 获取所有已分享的心愿
// @Tags wishes
// @Accept json
// @Produce json
// @Success 200 {array} model.SharedWish "分享的心愿列表"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /wishes/community [get]
func GetCommunityWishes(c *gin.Context) {
	var sharedWishes []model.SharedWish
	if err := database.GormDB.Find(&sharedWishes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取心愿社区失败"})
		return
	}

	c.JSON(http.StatusOK, sharedWishes)
}

// @Summary 获取随机心愿
// @Description 从心愿社区中随机获取一个心愿
// @Tags wishes
// @Accept json
// @Produce json
// @Success 200 {object} model.SharedWish "随机心愿"
// @Failure 404 {object} map[string]string "暂无共享心愿"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /wishes/random [get]
func GetRandomWish(c *gin.Context) {
	var count int64
	if err := database.GormDB.Model(&model.SharedWish{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取心愿数量失败"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "暂无共享心愿"})
		return
	}

	offset := rand.Int63n(count)
	var wish model.SharedWish
	if err := database.GormDB.Offset(int(offset)).First(&wish).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取随机心愿失败"})
		return
	}

	c.JSON(http.StatusOK, wish)
}

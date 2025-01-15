package v1

import (
	"net/http"
	"time"

	"github.com/PisaListBE/internal/model"
	"github.com/PisaListBE/pkg/database"
	"github.com/gin-gonic/gin"
)

// @title PisaList Task API
// @version 1.0
// @description This is the task management API for PisaList
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

type TaskRequest struct {
	Event           string `json:"event" binding:"required" example:"Buy groceries"`
	Description     string `json:"description" example:"Milk, eggs, bread"`
	IsCycle         bool   `json:"is_cycle" example:"false"`
	ImportanceLevel int    `json:"importance_level" binding:"min=0,max=5" example:"3"`
}

// @Summary 创建任务
// @Description 创建一个新的任务
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param task body TaskRequest true "任务信息"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	userID := c.GetUint("userID")
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := model.Task{
		UserID:          userID,
		Event:           req.Event,
		Description:     req.Description,
		IsCycle:         req.IsCycle,
		ImportanceLevel: req.ImportanceLevel,
	}

	if err := database.DB.Create(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建任务失败"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary 删除任务
// @Description 根据ID删除任务
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "任务ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var task model.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	if err := database.DB.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除任务失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// @Summary 更新任务
// @Description 更新任务信息
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "任务ID"
// @Param task body TaskRequest true "更新后的任务信息"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task model.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	updates := map[string]interface{}{
		"event":       req.Event,
		"description": req.Description,
		"is_cycle":    req.IsCycle,
		// "importance_level": req.ImportanceLevel,
	}

	if err := database.DB.Model(&task).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新任务失败"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary 完成任务
// @Description 标记任务为已完成
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "任务ID"
// @Success 200 {object} model.Task
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id}/complete [put]
func CompleteTask(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var task model.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	// 切换完成状态
	newCompleted := !task.Completed

	updates := map[string]interface{}{
		"completed": newCompleted,
	}

	if newCompleted {
		updates["completed_date"] = time.Now()
	} else {
		updates["completed_date"] = nil // 取消完成时清空完成日期
	}

	if err := database.DB.Model(&task).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新任务状态失败"})
		return
	}

	// 重新获取更新后的任务信息
	if err := database.DB.First(&task, task.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取更新后的任务失败"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary 更新任务优先级
// @Description 更新任务的重要性级别
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "任务ID"
// @Param importance body object{importance_level=int} true "新的优先级"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id}/importance [put]
func UpdateTaskImportance(c *gin.Context) {
	userID := c.GetUint("userID")
	taskID := c.Param("id")

	var req struct {
		ImportanceLevel int `json:"importance_level" binding:"required,min=0,max=5"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var task model.Task
	if err := database.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
		return
	}

	if err := database.DB.Model(&task).Update("importance_level", req.ImportanceLevel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新任务优先级失败"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary 获取任务时间线
// @Description 获取过去7天完成的任务
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Task
// @Failure 500 {object} map[string]string
// @Router /tasks/timeline [get]
func GetTaskTimeline(c *gin.Context) {
	userID := c.GetUint("userID")
	sevenDaysAgo := time.Now().AddDate(0, 0, -7) //计算出七天前的日期。这两个变量将用于后续的数据库查询。

	var tasks []model.Task
	if err := database.DB.Where("user_id = ? AND completed = ? AND completed_date >= ?",
		userID, true, sevenDaysAgo).
		Order("completed_date desc").
		Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取任务时间线失败"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// @Summary 获取今日任务
// @Description 获取今天需要完成的任务
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Task
// @Failure 500 {object} map[string]string
// @Router /tasks/today [get]
func GetTodayTasks(c *gin.Context) {
	userID := c.GetUint("userID")
	today := time.Now().Format("2006-01-02")

	var tasks []model.Task
	if err := database.DB.Where(
		"user_id = ? AND (completed = false OR (completed = true AND DATE(completed_date) = ?) OR (completed = true AND is_cycle = true))",
		userID,
		today,
	).Order("importance_level asc").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取今日任务失败"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// @Summary 获取用户任务
// @Description 获取用户所有任务
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.Task
// @Failure 500 {object} map[string]string
// @Router /tasks/user [get]
func GetUserTasks(c *gin.Context) {
	userID := c.GetUint("userID")

	var tasks []model.Task
	if err := database.DB.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户任务失败"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

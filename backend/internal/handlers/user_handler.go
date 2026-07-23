package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"inventory-management/backend/internal/models"
)

type UserHandler struct {
	db *gorm.DB
}

type UserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	RoleID      uint   `json:"role_id" binding:"required"`
}

type UserResponse struct {
	ID          uint     `json:"id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Status      string   `json:"status"`
	RoleID      uint     `json:"role_id"`
	Role        RoleInfo `json:"role"`
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := h.db.Preload("Role").Order("id ASC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取用户失败"})
		return
	}

	items := make([]UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, buildUserResponse(user))
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名和角色不能为空"})
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "新建用户必须设置密码"})
		return
	}

	user, err := h.buildUserFromRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "密码加密失败"})
		return
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "创建用户失败，请检查用户名是否重复"})
		return
	}

	if err := h.db.Preload("Role").First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取用户失败"})
		return
	}

	c.JSON(http.StatusCreated, buildUserResponse(user))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名和角色不能为空"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		return
	}

	user.Username = req.Username
	user.DisplayName = req.DisplayName
	user.RoleID = req.RoleID
	user.Status = normalizeUserStatus(req.Status)

	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "密码加密失败"})
			return
		}
		user.PasswordHash = string(hash)
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "更新用户失败，请检查用户名是否重复"})
		return
	}

	if err := h.db.Preload("Role").First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取用户失败"})
		return
	}

	c.JSON(http.StatusOK, buildUserResponse(user))
}

func (h *UserHandler) ToggleUserStatus(c *gin.Context) {
	userID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户状态不能为空"})
		return
	}

	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "用户不存在"})
		return
	}

	user.Status = normalizeUserStatus(req.Status)
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "更新用户状态失败"})
		return
	}

	if err := h.db.Preload("Role").First(&user, user.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取用户失败"})
		return
	}

	c.JSON(http.StatusOK, buildUserResponse(user))
}

func (h *UserHandler) buildUserFromRequest(req UserRequest) (models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Username:     req.Username,
		PasswordHash: string(hash),
		DisplayName:  req.DisplayName,
		Status:       normalizeUserStatus(req.Status),
		RoleID:       req.RoleID,
	}, nil
}

func buildUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Status:      user.Status,
		RoleID:      user.RoleID,
		Role: RoleInfo{
			ID:   user.Role.ID,
			Code: user.Role.Code,
			Name: user.Role.Name,
		},
	}
}

func normalizeUserStatus(status string) string {
	if status == models.UserStatusDisabled {
		return models.UserStatusDisabled
	}
	return models.UserStatusEnabled
}

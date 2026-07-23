package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"inventory-management/backend/internal/config"
	"inventory-management/backend/internal/models"
	"inventory-management/backend/internal/services"
)

type AuthHandler struct {
	db  *gorm.DB
	cfg config.Config
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserResponse struct {
	ID          uint     `json:"id"`
	Username    string   `json:"username"`
	DisplayName string   `json:"display_name"`
	Role        RoleInfo `json:"role"`
	Permissions []string `json:"permissions"`
}

type LoginResponse struct {
	Token     string            `json:"token"`
	ExpiresAt int64             `json:"expires_at"`
	User      LoginUserResponse `json:"user"`
}

type RoleInfo struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

func NewAuthHandler(db *gorm.DB, cfg config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "用户名和密码不能为空"})
		return
	}

	var user models.User
	err := h.db.Preload("Role").
		Where("username = ?", req.Username).
		First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		return
	}

	if user.Status != models.UserStatusEnabled {
		c.JSON(http.StatusForbidden, gin.H{"message": "用户已被禁用"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户名或密码错误"})
		return
	}

	permissions, err := h.getPermissionCodes(user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取权限失败"})
		return
	}

	token, expiresAt, err := services.GenerateToken(h.cfg.JWTSecret, user.ID, user.Username, user.Role.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成登录凭证失败"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      buildLoginUser(user, permissions),
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	value, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "未登录"})
		return
	}

	user, ok := value.(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "登录状态异常"})
		return
	}

	permissions, err := h.getPermissionCodes(user.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取权限失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": buildLoginUser(user, permissions)})
}

func (h *AuthHandler) getPermissionCodes(roleID uint) ([]string, error) {
	var permissions []models.Permission
	err := h.db.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.module ASC, permissions.code ASC").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	codes := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		codes = append(codes, permission.Code)
	}

	return codes, nil
}

func buildLoginUser(user models.User, permissions []string) LoginUserResponse {
	return LoginUserResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role: RoleInfo{
			ID:   user.Role.ID,
			Code: user.Role.Code,
			Name: user.Role.Name,
		},
		Permissions: permissions,
	}
}

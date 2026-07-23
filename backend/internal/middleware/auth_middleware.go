package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"inventory-management/backend/internal/config"
	"inventory-management/backend/internal/models"
	"inventory-management/backend/internal/services"
)

func AuthRequired(db *gorm.DB, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "未登录"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录凭证格式错误"})
			return
		}

		claims, err := services.ParseToken(cfg.JWTSecret, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录已过期，请重新登录"})
			return
		}

		var user models.User
		err = db.Preload("Role").
			Where("id = ?", claims.UserID).
			First(&user).Error
		if err != nil || user.Status != models.UserStatusEnabled {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "用户不存在或已禁用"})
			return
		}

		c.Set("current_user", user)
		c.Next()
	}
}

func RequirePermission(db *gorm.DB, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get("current_user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "未登录"})
			return
		}

		user, ok := value.(models.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "登录状态异常"})
			return
		}

		var count int64
		err := db.Model(&models.RolePermission{}).
			Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
			Where("role_permissions.role_id = ? AND permissions.code = ?", user.RoleID, permissionCode).
			Count(&count).Error
		if err != nil || count == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "无权限"})
			return
		}

		c.Next()
	}
}

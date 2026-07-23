package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"inventory-management/backend/internal/models"
)

type RoleHandler struct {
	db *gorm.DB
}

type RoleRequest struct {
	Code            string   `json:"code" binding:"required"`
	Name            string   `json:"name" binding:"required"`
	Description     string   `json:"description"`
	PermissionCodes []string `json:"permission_codes"`
}

type RoleResponse struct {
	ID              uint     `json:"id"`
	Code            string   `json:"code"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	IsSystem        bool     `json:"is_system"`
	PermissionCodes []string `json:"permission_codes"`
}

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Module      string `json:"module"`
	Description string `json:"description"`
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{db: db}
}

func (h *RoleHandler) ListRoles(c *gin.Context) {
	var roles []models.Role
	if err := h.db.Order("id ASC").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取角色失败"})
		return
	}

	roleIDs := make([]uint, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	permissionMap, err := h.getRolePermissionMap(roleIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取角色权限失败"})
		return
	}

	items := make([]RoleResponse, 0, len(roles))
	for _, role := range roles {
		items = append(items, buildRoleResponse(role, permissionMap[role.ID]))
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *RoleHandler) ListPermissions(c *gin.Context) {
	var permissions []models.Permission
	if err := h.db.Order("module ASC, code ASC").Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "读取权限失败"})
		return
	}

	items := make([]PermissionResponse, 0, len(permissions))
	for _, permission := range permissions {
		items = append(items, PermissionResponse{
			ID:          permission.ID,
			Code:        permission.Code,
			Name:        permission.Name,
			Module:      permission.Module,
			Description: permission.Description,
		})
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "角色标识和角色名称不能为空"})
		return
	}

	var role models.Role
	err := h.db.Transaction(func(tx *gorm.DB) error {
		role = models.Role{
			Code:        req.Code,
			Name:        req.Name,
			Description: req.Description,
			IsSystem:    false,
		}
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		return h.replaceRolePermissions(tx, role.ID, req.PermissionCodes)
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "创建角色失败，请检查角色标识是否重复"})
		return
	}

	c.JSON(http.StatusCreated, buildRoleResponse(role, req.PermissionCodes))
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	roleID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var req RoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "角色标识和角色名称不能为空"})
		return
	}

	var role models.Role
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		role.Code = req.Code
		role.Name = req.Name
		role.Description = req.Description
		if err := tx.Save(&role).Error; err != nil {
			return err
		}

		return h.replaceRolePermissions(tx, role.ID, req.PermissionCodes)
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "更新角色失败，请检查角色是否存在或标识是否重复"})
		return
	}

	c.JSON(http.StatusOK, buildRoleResponse(role, req.PermissionCodes))
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	roleID, ok := parseIDParam(c, "id")
	if !ok {
		return
	}

	var role models.Role
	if err := h.db.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "角色不存在"})
		return
	}
	if role.IsSystem {
		c.JSON(http.StatusBadRequest, gin.H{"message": "系统内置角色不能删除"})
		return
	}

	var userCount int64
	if err := h.db.Model(&models.User{}).Where("role_id = ?", role.ID).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "检查角色使用状态失败"})
		return
	}
	if userCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "该角色已有用户使用，不能删除"})
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}
		return tx.Delete(&role).Error
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "删除角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func (h *RoleHandler) getRolePermissionMap(roleIDs []uint) (map[uint][]string, error) {
	result := make(map[uint][]string)
	if len(roleIDs) == 0 {
		return result, nil
	}

	type row struct {
		RoleID         uint
		PermissionCode string
	}
	var rows []row
	err := h.db.Model(&models.RolePermission{}).
		Select("role_permissions.role_id, permissions.code AS permission_code").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Order("permissions.module ASC, permissions.code ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, item := range rows {
		result[item.RoleID] = append(result[item.RoleID], item.PermissionCode)
	}

	return result, nil
}

func (h *RoleHandler) replaceRolePermissions(tx *gorm.DB, roleID uint, permissionCodes []string) error {
	if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}
	if len(permissionCodes) == 0 {
		return nil
	}

	var permissions []models.Permission
	if err := tx.Where("code IN ?", permissionCodes).Find(&permissions).Error; err != nil {
		return err
	}
	if len(permissions) != len(permissionCodes) {
		return gorm.ErrRecordNotFound
	}

	for _, permission := range permissions {
		rolePermission := models.RolePermission{
			RoleID:       roleID,
			PermissionID: permission.ID,
		}
		if err := tx.Create(&rolePermission).Error; err != nil {
			return err
		}
	}

	return nil
}

func buildRoleResponse(role models.Role, permissionCodes []string) RoleResponse {
	return RoleResponse{
		ID:              role.ID,
		Code:            role.Code,
		Name:            role.Name,
		Description:     role.Description,
		IsSystem:        role.IsSystem,
		PermissionCodes: permissionCodes,
	}
}

func parseIDParam(c *gin.Context, key string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil || value == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "参数错误"})
		return 0, false
	}

	return uint(value), true
}

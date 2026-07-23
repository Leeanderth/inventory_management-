package main

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"inventory-management/backend/internal/config"
	"inventory-management/backend/internal/database"
	"inventory-management/backend/internal/models"
)

type seedUser struct {
	Username    string
	DisplayName string
	RoleCode    string
	PasswordEnv string
	DefaultPass string
}

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	if err := seedPermissions(db); err != nil {
		log.Fatalf("seed permissions: %v", err)
	}
	if err := seedRoles(db); err != nil {
		log.Fatalf("seed roles: %v", err)
	}
	if err := seedRolePermissions(db); err != nil {
		log.Fatalf("seed role permissions: %v", err)
	}
	if err := seedUsers(db); err != nil {
		log.Fatalf("seed users: %v", err)
	}

	log.Println("seed completed")
}

func seedPermissions(db *gorm.DB) error {
	permissions := []models.Permission{
		{Code: "user:view", Name: "查看用户列表", Module: "user"},
		{Code: "user:create", Name: "创建用户", Module: "user"},
		{Code: "user:update", Name: "编辑用户", Module: "user"},
		{Code: "user:disable", Name: "启用或禁用用户", Module: "user"},
		{Code: "role:view", Name: "查看角色列表", Module: "role"},
		{Code: "role:create", Name: "创建角色", Module: "role"},
		{Code: "role:update", Name: "编辑角色", Module: "role"},
		{Code: "role:delete", Name: "删除角色", Module: "role"},
		{Code: "stock:view", Name: "查看库存", Module: "stock"},
		{Code: "stock:create", Name: "新增商品", Module: "stock"},
		{Code: "stock:update", Name: "编辑商品", Module: "stock"},
		{Code: "stock:delete", Name: "删除商品", Module: "stock"},
	}

	for _, permission := range permissions {
		if err := db.Where("code = ?", permission.Code).FirstOrCreate(&permission).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedRoles(db *gorm.DB) error {
	roles := []models.Role{
		{Code: "super_admin", Name: "超级管理员", Description: "拥有系统全部权限", IsSystem: true},
		{Code: "stock_manager", Name: "库存管理员", Description: "可以查看、新增、编辑、删除库存", IsSystem: true},
		{Code: "stock_viewer", Name: "库存查看员", Description: "只能查看库存", IsSystem: true},
	}

	for _, role := range roles {
		if err := db.Where("code = ?", role.Code).FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedRolePermissions(db *gorm.DB) error {
	rolePermissions := map[string][]string{
		"super_admin": {
			"user:view", "user:create", "user:update", "user:disable",
			"role:view", "role:create", "role:update", "role:delete",
			"stock:view", "stock:create", "stock:update", "stock:delete",
		},
		"stock_manager": {"stock:view", "stock:create", "stock:update", "stock:delete"},
		"stock_viewer":  {"stock:view"},
	}

	for roleCode, permissionCodes := range rolePermissions {
		var role models.Role
		if err := db.Where("code = ?", roleCode).First(&role).Error; err != nil {
			return err
		}

		for _, permissionCode := range permissionCodes {
			var permission models.Permission
			if err := db.Where("code = ?", permissionCode).First(&permission).Error; err != nil {
				return err
			}

			rolePermission := models.RolePermission{
				RoleID:       role.ID,
				PermissionID: permission.ID,
			}
			if err := db.Where("role_id = ? AND permission_id = ?", role.ID, permission.ID).
				FirstOrCreate(&rolePermission).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func seedUsers(db *gorm.DB) error {
	users := []seedUser{
		{Username: "root", DisplayName: "超级管理员", RoleCode: "super_admin", PasswordEnv: "ROOT_PASSWORD", DefaultPass: "Root@123456"},
		{Username: "manager", DisplayName: "库存管理员", RoleCode: "stock_manager", PasswordEnv: "MANAGER_PASSWORD", DefaultPass: "Manager@123456"},
		{Username: "viewer", DisplayName: "库存查看员", RoleCode: "stock_viewer", PasswordEnv: "VIEWER_PASSWORD", DefaultPass: "Viewer@123456"},
	}

	for _, item := range users {
		var role models.Role
		if err := db.Where("code = ?", item.RoleCode).First(&role).Error; err != nil {
			return err
		}

		password := os.Getenv(item.PasswordEnv)
		if password == "" {
			password = item.DefaultPass
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := models.User{
			Username:     item.Username,
			PasswordHash: string(hash),
			DisplayName:  item.DisplayName,
			Status:       models.UserStatusEnabled,
			RoleID:       role.ID,
		}

		if err := db.Where("username = ?", item.Username).FirstOrCreate(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

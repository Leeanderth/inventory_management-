package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"inventory-management/backend/internal/config"
	"inventory-management/backend/internal/handlers"
	"inventory-management/backend/internal/middleware"
)

func Setup(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		healthHandler := handlers.NewHealthHandler(db)
		api.GET("/health", healthHandler.Check)

		authHandler := handlers.NewAuthHandler(db, cfg)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthRequired(db, cfg), authHandler.Me)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthRequired(db, cfg))
		{
			roleHandler := handlers.NewRoleHandler(db)
			userHandler := handlers.NewUserHandler(db)
			productHandler := handlers.NewProductHandler(db)

			protected.GET("/permissions", middleware.RequirePermission(db, "role:view"), roleHandler.ListPermissions)
			protected.GET("/roles", middleware.RequirePermission(db, "role:view"), roleHandler.ListRoles)
			protected.POST("/roles", middleware.RequirePermission(db, "role:create"), roleHandler.CreateRole)
			protected.PUT("/roles/:id", middleware.RequirePermission(db, "role:update"), roleHandler.UpdateRole)
			protected.DELETE("/roles/:id", middleware.RequirePermission(db, "role:delete"), roleHandler.DeleteRole)

			protected.GET("/users", middleware.RequirePermission(db, "user:view"), userHandler.ListUsers)
			protected.POST("/users", middleware.RequirePermission(db, "user:create"), userHandler.CreateUser)
			protected.PUT("/users/:id", middleware.RequirePermission(db, "user:update"), userHandler.UpdateUser)
			protected.PATCH("/users/:id/status", middleware.RequirePermission(db, "user:disable"), userHandler.ToggleUserStatus)

			protected.GET("/products", middleware.RequirePermission(db, "stock:view"), productHandler.ListProducts)
			protected.POST("/products", middleware.RequirePermission(db, "stock:create"), productHandler.CreateProduct)
			protected.PUT("/products/:id", middleware.RequirePermission(db, "stock:update"), productHandler.UpdateProduct)
			protected.DELETE("/products/:id", middleware.RequirePermission(db, "stock:delete"), productHandler.DeleteProduct)
			protected.GET("/products/:id/movements", middleware.RequirePermission(db, "stock:view"), productHandler.ListStockMovements)
		}
	}

	return router
}

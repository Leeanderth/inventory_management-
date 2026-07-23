package main

import (
	"log"

	"inventory-management/backend/internal/config"
	"inventory-management/backend/internal/database"
	"inventory-management/backend/internal/routes"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	router := routes.Setup(db, cfg)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("start server: %v", err)
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/config"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/database"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config Error")
	}

	client, db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("db error")
	}

	defer func() {
		if err := database.Disconnect(client); err != nil {
			log.Printf("momgo disconnect error: %v", err)
		}
	}()

	router := handler.NewRouter(db)
	addr := fmt.Sprintf(":%s", cfg.ServerPort)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server Failed")
	}
}

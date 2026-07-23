package main

import (
	"fmt"
	"log"

	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/config"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/database"
	"github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		//log.Fatalf("Config Error")
		log.Fatalf("Config Error: %v", err)
	}

	client, db, err := database.Connect(cfg)
	if err != nil {
		//log.Fatalf("db error")
		log.Fatalf("DB Error: %v", err)
	}

	defer func() {
		if err := database.Disconnect(client); err != nil {
			log.Printf("mongo disconnect error: %v", err)
		}
	}()

	router := router.NewRouter(db, cfg)

	addr := fmt.Sprintf(":%s", cfg.ServerPort)

	log.Printf("Server listening on http://localhost%s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server Failed")
	}

}

package main

import (
	"log"

	"github.com/ashil-poojary/gostruct/internal/config"
	"github.com/ashil-poojary/gostruct/internal/db"
	"github.com/ashil-poojary/gostruct/internal/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// init db
	if err := db.InitDB1(cfg.Databases["db1"]); err != nil {
		log.Fatalf("Error initializing DB1: %v", err)
	}
	if err := db.InitDB2(cfg.Databases["db2"]); err != nil {
		log.Fatalf("Error initializing DB2: %v", err)
	}
	if err := db.InitDB3(cfg.Databases["db3"]); err != nil {
		log.Fatalf("Error initializing DB3: %v", err)
	}

	r := routes.SetupRouter(cfg)
	log.Printf("Starting server on :%s", cfg.ServerPort)
	log.Fatal(r.Run(":" + cfg.ServerPort))
}

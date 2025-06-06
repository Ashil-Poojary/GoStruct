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
	if err := db.InitDefaultDB(cfg.DefaultDB); err != nil {
		log.Fatalf("Error initializing DefaultDB: %v", err)
	}
	if err := db.InitReplicaDB(cfg.ReplicaDB); err != nil {
		log.Fatalf("Error initializing ReplicaDB: %v", err)
	}
	if err := db.InitProdDB(cfg.ProdDB); err != nil {
		log.Fatalf("Error initializing ProdDB: %v", err)
	}

	r := routes.SetupRouter(cfg)
	log.Printf("Starting server on :%s", cfg.ServerPort)
	log.Fatal(r.Run(":" + cfg.ServerPort))
}

package db

import (
	"fmt"
	"log"

	"github.com/ashil-poojary/gostruct/internal/config"
	"github.com/ashil-poojary/gostruct/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB1 *gorm.DB
	DB2 *gorm.DB
	DB3 *gorm.DB
)

func connectDB(cfg config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	// Open a new database connection

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}

func InitDB1(cfg config.DBConfig) error {
	db, err := connectDB(cfg)
	if err != nil {
		return err
	}
	DB1 = db

	if cfg.Migration {
		log.Println("Migrating DB1 models...")
		if err := db.AutoMigrate(
			&models.User{},
			&models.RefreshToken{},

			&models.Role{},
			&models.Order{},
			&models.OrderReturn{},
		); err != nil {
			return err
		}
	}

	return nil
}

func InitDB2(cfg config.DBConfig) error {
	db, err := connectDB(cfg)
	if err != nil {
		return err
	}
	DB2 = db

	if cfg.Migration {
		log.Println("Migrating DB2 models...")
		if err := db.AutoMigrate(
			&models.RefreshToken{},
			&models.User{},
			&models.Role{},
			&models.Order{},
			&models.OrderReturn{},

			// add models for DB2 here
		); err != nil {
			return err
		}
	}

	return nil
}

func InitDB3(cfg config.DBConfig) error {
	db, err := connectDB(cfg)
	if err != nil {
		return err
	}
	DB3 = db

	if cfg.Migration {
		log.Println("Migrating DB3 models...")
		if err := db.AutoMigrate(
			&models.User{},
			&models.RefreshToken{},

			&models.Role{},
			&models.Order{},
			&models.OrderReturn{},

			// add models for DB3 here
		); err != nil {
			return err
		}
	}

	return nil
}

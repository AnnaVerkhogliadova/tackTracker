package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connect() *pgxpool.Pool {
	dsn := "host=localhost user=user password=1234 dbname=tasks port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	log.Println("Database connected (GORM)")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Could not parse pgxpool config: %v", err)
	}

	rwdb, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Could not create pgxpool: %v", err)
	}

	return rwdb
}

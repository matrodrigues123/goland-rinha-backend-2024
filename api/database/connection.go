package database

import (
	"fmt"

	"github.com/matrodrigues123/rinha-2024-go-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "db"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbName   = "rinha-db"
)



func Connection() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) 
	if err != nil {
		panic(err.Error())
	}

	migrator := db.Migrator()
	if !migrator.HasTable(&models.Account{}) || !migrator.HasTable(&models.Transaction{}) {
		err = db.AutoMigrate(&models.Account{}, &models.Transaction{})
		if err != nil {
			panic(fmt.Sprintf("failed to automigrate: %v", err))
		}

		accounts := []*models.Account{
			{Limit: 1e5, Total: 0},
			{Limit: 8e4, Total: 0},
			{Limit: 1e6, Total: 0},
			{Limit: 1e7, Total: 0},
			{Limit: 5e5, Total: 0},
		}
		db.Create(accounts)
	}
	
	return db

}

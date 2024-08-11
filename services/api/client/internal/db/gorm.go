package db

import (
	"client/internal/db/entities"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	ComputeId string
	Password  string
	DB_Name   string
}

func NewConfig(v *viper.Viper) *Config {
	return &Config{
		ComputeId: v.GetString("db.computeId"),
		Password:  v.GetString("db.password"),
		DB_Name:   v.GetString("db.dbName"),
	}
}

func (cf *Config) PostgresConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgresql://virtual_banking_db_owner:%s@%s.ap-southeast-1.aws.neon.tech/%s?sslmode=require", cf.Password, cf.ComputeId, cf.DB_Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Connected to database")

	cf.Migrate(db)

	return db, nil
}

func (cf *Config) Migrate(db *gorm.DB) {
	db.AutoMigrate(&entities.Client{})
	db.AutoMigrate(&entities.Account{})
	db.AutoMigrate(&entities.Transaction{})
	db.AutoMigrate(&entities.Payment{})
}

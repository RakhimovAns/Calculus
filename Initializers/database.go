package Initializers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

type Expression struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	Expression   string
	AddTime      time.Time `gorm:"type:timestamp"`
	SubTime      time.Time `gorm:"type:timestamp"`
	MultiplyTime time.Time `gorm:"type:timestamp"`
	DivideTime   time.Time `gorm:"type:timestamp"`
	Created      time.Time `gorm:"type:timestamp"`
	Finished     time.Time `gorm:"type:timestamp"`
}

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	dsn = "host=localhost user=postgres password=postgres dbname=yandex port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error with connection to database")
	}
}
func CreateTable() {
	err := DB.AutoMigrate(&Expression{})
	if err != nil {
		log.Fatal("failed to migrate")
	}
}

package database

import (
	"fmt"
	"log"
	"vix-btpns/helpers"
	"vix-btpns/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var  (
	db *gorm.DB
	err error
)

func InitDB() {
	username := helpers.GetEnv("MYSQL_USER_LOCAL")
	password := helpers.GetEnv("MYSQL_PASSWORD_LOCAL")
	host := helpers.GetEnv("MYSQL_HOST_LOCAL")
	port := helpers.GetEnv("MYSQL_PORT_LOCAL")
	dbName := helpers.GetEnv("MYSQL_DB_NAME_LOCAL")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbName)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	doMigration()
}

func GetDB() *gorm.DB {
	return db
}

func doMigration() {
	if err := db.AutoMigrate(
		&models.User{},
		&models.UserPhoto{},
	); err != nil {
		fmt.Println("Migration failed!")
		log.Fatal(err.Error())
	}
}

package db

import (
	"fmt"
	"github.com/hewo233/hdu-dx2/models"
	"github.com/hewo233/hdu-dx2/shared/consts"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func UpdateDB() {
	err := DB.Table(consts.UserTable).AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}
	err = DB.Table(consts.FamilyTable).AutoMigrate(&models.Family{})
	if err != nil {
		log.Fatal(err)
	}
	err = DB.Table(consts.FamilyUserTable).AutoMigrate(&models.FamilyUser{})
	if err != nil {
		log.Fatal(err)
	}
	err = DB.Table(consts.BillTable).AutoMigrate(&models.Bill{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("\033[32mAutoMigrate success\033[0m")
}

func ConnectDB() {

	if err := godotenv.Load(consts.DBEnvFile); err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, port)

	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
}

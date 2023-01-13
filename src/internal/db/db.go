package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"package/main/internal/config"
	"package/main/internal/models"
)

var conf *config.Config
var DB *gorm.DB
var err error

func init() {
	conf = config.Cfg
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN: "host=" + conf.DB_HOST +
			" user=" + conf.DB_USER +
			" password=" + conf.DB_PASS +
			" dbname=" + conf.DB_DATABASE +
			" port=" + conf.DB_PORT +
			" sslmode=disable", // TimeZone=
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		fmt.Printf("Something went wrong while connecting to database: %s", err)
		os.Exit(1)
	}

	DB.AutoMigrate(&models.User{})
}

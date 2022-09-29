package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Instance *gorm.DB

func Init() {
	var err error
	Instance, err = gorm.Open(sqlite.Open(os.Getenv("DB")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

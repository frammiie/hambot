package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Instance *gorm.DB

func Init() {
	var err error
	
	Instance, err = gorm.Open(sqlite.Open(os.Getenv("DB")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Minute,
				IgnoreRecordNotFoundError: true,
			},
		),
	})

	if err != nil {
		log.Fatal(err)
	}
}

package usm

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var client *Client

func TestMain(m *testing.M) {
	user := ""
	password := ""
	host := ""
	port := ""
	database := ""

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		database)

	var err error
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:                   false,
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Warn, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,        // Don't include params in the SQL log
				Colorful:                  false,       // Disable color
			},
		),
	})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	client, err = New(&Options{
		Store: db,
		CasbinOptions: &CasbinOptions{
			Store:      db,
			PolicyPath: "conf/casbin.conf",
		},
	})

	m.Run()
}

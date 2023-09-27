package config

import (
	app "com.nandini.icici/config"
	"fmt"
	"github.com/micro/micro/v3/service/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var stmt = "create table if not exists %v(id text not null, title text,description text, primary key(id));"
var tableName = "event"

func GetClient() *gorm.DB {
	db_url := app.GetVal("GO_MICRO_DB_URL")
	db, err := gorm.Open(postgres.Open(db_url), &gorm.Config{})
	if err != nil {
		logger.Errorf("failed to connect database")
		return nil
	}
	logger.Infof("Database connected")
	db.Exec(fmt.Sprintf(stmt, tableName))
	return db
}

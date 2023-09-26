package migrate

import (
	"github.com/micro/micro/v3/service/logger"
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"strings"
	app "com.nandini.icici/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func generateUrl(databaseUrl, defaultDatabaseName, newDatabaseName string) string {
	return func() string {
		if lastIndex := strings.LastIndex(databaseUrl, defaultDatabaseName); lastIndex != -1 {
			return databaseUrl[:lastIndex] + newDatabaseName + databaseUrl[lastIndex+len(defaultDatabaseName):]
		}
		return databaseUrl
	}()
}


func MigrateAndCreateDatabase() {
	dsn := app.GetVal("GO_MICRO_DB_URL")+"?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Errorf(err.Error())
	}
	defer db.Close()
	newUrl :=dsn
	env := app.GetVal("GO_MICRO_ACTIVE_PROFILE")
	if(env=="dev"){
		dbName := "icici"
		var exists bool
		query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
		err = db.QueryRow(query).Scan(&exists)
		if err != nil {
			logger.Errorf("Failed to check if the database exists:", err)
		}

		if !exists {
			createDBSQL := fmt.Sprintf("CREATE DATABASE %s", dbName)
			_, err := db.Exec(createDBSQL)
			if err != nil {
				logger.Errorf("Failed to create the database:", err)
			}
			logger.Infof("Database created successfully!")
		} else {
			logger.Infof("Database already exists.")
		}
		newUrl =generateUrl(dsn,"postgres",dbName)
		app.SetVal("GO_MICRO_DB_URL",newUrl)
	}
	migrator, err := migrate.New(
		"file://migrate/migration",
		newUrl)
	if err != nil {
			logger.Errorf(err.Error())
	}
	if err := migrator.Up(); err != nil {
			logger.Errorf(err.Error())
	}
}
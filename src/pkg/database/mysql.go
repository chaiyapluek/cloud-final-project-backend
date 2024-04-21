package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"dev.chaiyapluek.cloud.final.backend/src/config"
	_ "github.com/go-sql-driver/mysql"
)

func Connect(cfg *config.DBConfig) *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.DBName))
	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

package initializer

import (
	"auth-service/helper"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgressDB() *sql.DB {
	dsn := helper.PostgresDSNBuilder("postgres", "admin", "localhost", "5432", "catwitch-dev");
	
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)
	db.SetConnMaxIdleTime(1 * time.Hour)
	db.SetConnMaxLifetime(3 * time.Hour)

	return db
}
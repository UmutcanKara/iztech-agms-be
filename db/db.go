package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, database)
	//pgSetup := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, database, password)
	pgSetup := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)
	//pgSetup := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)

	db, err := sql.Open("postgres", pgSetup)
	if err != nil {
		return nil, err
	}
	return &Database{db: db}, nil
}
func (d *Database) Close() {
	err := d.db.Close()
	if err != nil {
		return
	}
}
func (d *Database) GetDB() *sql.DB {
	return d.db
}

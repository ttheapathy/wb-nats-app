package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDatabase() *sqlx.DB {

	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL"),
	)

	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createSchema(db *sqlx.DB) {
	schema := `CREATE TABLE IF NOT EXISTS messages (
		id         serial  not null unique,
		text       varchar(120) NOT NULL
	);`

	db.MustExec(schema)
}

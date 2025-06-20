package database

import (
	"database/sql"
	"log"

	"github.com/rubenv/sql-migrate"
)

func GetDb() *sql.DB {
	migrations := &migrate.FileMigrationSource{Dir: "../../migrations"}
	db, err := sql.Open("sqlite3", "../../dota-pro-games.db?foreign_keys=true")
	if err != nil {
		log.Fatal(err)
	}
	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Applied %d migrations!\n", n)
	return db
}

package database

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rubenv/sql-migrate"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func GetDb() *sql.DB {
	useTurso := flag.Bool("turso", false, "Uses remote turso db")
	flag.Parse()
	if *useTurso == true {
		return getTursoDb()
	}
	return getSqliteDb()
}

func getSqliteDb() *sql.DB {
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

func getTursoDb() *sql.DB {
	var databaseUrl = os.Getenv("TURSO_DATABASE_URL")
	if databaseUrl == "" {
		log.Fatalln("TURSO_DATABASE_URL is not set")
	}
	var authToken = os.Getenv("TURSO_AUTH_TOKEN")
	if authToken == "" {
		log.Fatalln("TURSO_AUTH_TOKEN is not set")
	}

	url := fmt.Sprintf("%s?authToken=%s", databaseUrl, authToken)
	db, err := sql.Open("libsql", url)

	migrations := &migrate.FileMigrationSource{Dir: "../../migrations"}
	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Applied %d migrations!\n", n)
	if err != nil {
		log.Fatalf("Failed to open turso db %s: %s \n", url, err)
	}
	return db

}

package main

import (
	"context"
	"database/sql"
	"dota-pro-db/database"
	"dota-pro-db/stratz"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	ctx := context.Background()
	client := stratz.GetClient()
	heroesRes, err := stratz.GetHeroes(client, ctx)
	if err != nil {
		log.Fatal(err)
	}

	heroes := heroesRes.Constants.Heroes

	itemsRes, err := stratz.GetItems(client, ctx)
	if err != nil {
		log.Fatal(err)
	}

	items := itemsRes.Constants.Items

	gameVersionsRes, err := stratz.GetGameVersions(client, ctx)
	if err != nil {
		log.Fatal(err)
	}
	gameVersions := gameVersionsRes.Constants.GameVersions

	db := database.GetDb()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	insertItems(db, items)
	insertHeroes(db, heroes)
	insertGameVersions(db, gameVersions)

}

func insertHeroes(db *sql.DB, heroes []stratz.Hero) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO heroes(id, display_name, short_name, primary_attribute) values(?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, hero := range heroes {
		_, err = stmt.Exec(hero.Id, hero.DisplayName, hero.ShortName, hero.Stats.PrimaryAttributeEnum)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted Heroes")
}
func insertItems(db *sql.DB, items []stratz.Item) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO items(id, display_name, short_name) values(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, item := range items {
		_, err = stmt.Exec(item.Id, item.DisplayName, item.ShortName)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted Items")
}

func insertGameVersions(db *sql.DB, versions []stratz.GameVersion) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT OR REPLACE INTO game_versions(id, name, asOfDateTime) VALUES(?,?,?)")
	defer stmt.Close()
	for _, version := range versions {
		_, err := stmt.Exec(version.Id, version.Name, version.AsOfDateTime)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted game versions")
}

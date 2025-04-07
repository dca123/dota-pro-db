package main

import (
	"context"
	"database/sql"
	"dota-pro-db/database"
	"dota-pro-db/stratz"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Khan/genqlient/graphql"
	_ "github.com/mattn/go-sqlite3"
)

var API_KEY = os.Getenv("STRATZ_API_KEY")
var URL = "https://api.stratz.com/graphql"

type authedTransport struct {
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "STRATZ_API")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", API_KEY))
	return t.wrapped.RoundTrip(req)
}
func main() {
	if API_KEY == "" {
		log.Fatalln("STRATZ_API_KEY is not set")
	}
	ctx := context.Background()
	client := graphql.NewClient(URL, &http.Client{
		Transport: &authedTransport{
			wrapped: http.DefaultTransport,
		},
	})

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

	db := database.GetDb()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	insertItems(db, items)
	insertHeroes(db, heroes)

}

func insertHeroes(db *sql.DB, heroes []stratz.Hero) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO heroes(id, display_name, short_name) values(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, hero := range heroes {
		_, err = stmt.Exec(hero.Id, hero.DisplayName, hero.ShortName)
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

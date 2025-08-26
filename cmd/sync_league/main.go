package main

import (
	"context"
	"dota-pro-db/create_league"
	"dota-pro-db/database"
	"dota-pro-db/stratz"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()
	matchId, err := strconv.Atoi(os.Args[1])
	liquipediaName := os.Args[2]
	if err != nil {
		log.Fatal(err)
	}
	if liquipediaName == "" {
		log.Fatal("Liquipedia name cannot be empty")
	}

	client := stratz.GetClient()

	fmt.Printf("Importing data for league: %d\n", matchId)
	league, err := stratz.GetLeagueMatches(client, ctx, matchId)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("League: ", league.DisplayName, "No. of matches: ", len(league.Matches))

	db := database.GetDb()
	defer db.Close()
	err = create_league.CreateLeague(db, &league, liquipediaName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully imported league")

}

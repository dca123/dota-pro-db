package main

import (
	"context"
	"dota-pro-db/create_league"
	"dota-pro-db/database"
	"dota-pro-db/liquipedia"
	"dota-pro-db/stratz"
	"fmt"
	"log"
)

type League = struct {
	Id int
}

func main() {
	ctx := context.Background()
	client := stratz.GetClient()

	leagues, err := liquipedia.GetLeagueIds()
	fmt.Printf("Found %d leagues\n", len(leagues))
	if err != nil {
		log.Fatal(err)

	}

	db := database.GetDb()
	defer db.Close()

	fmt.Println("--League Info--")
	fmt.Printf("Importing %d leagues\n", len(leagues))
	for _, league := range leagues {
		fmt.Println("📥 Importing league", league.ID)
		leagueWithMatches, err := stratz.GetLeagueMatches(client, ctx, league.ID)
		if err != nil {
			log.Fatal(err)
		}

		err = create_league.CreateLeague(db, &leagueWithMatches, league.Name)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Imported all leagues")

}

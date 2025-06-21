package main

import (
	"context"
	"database/sql"
	"dota-pro-db/create_league"
	"dota-pro-db/database"
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

	const startDateTime = 1739644200
	const minPrizePool = 1_000_000
	fmt.Printf("Getting leagues > $%d from %d onwards\n", minPrizePool, startDateTime)
	//gets leagues that have ended by today
	leagues, err := stratz.GetLeagues(client, ctx, startDateTime, minPrizePool)
	fmt.Printf("Found %d leagues\n", len(leagues))
	if err != nil {
		log.Fatal(err)

	}

	db := database.GetDb()
	defer db.Close()

	fmt.Println("--League Info--")
	var leaguesToImport []League
	for _, league := range leagues {
		var scannedLeague int
		err = db.QueryRow("SELECT id FROM leagues WHERE id = ?", league.Id).Scan(&scannedLeague)
		if err != nil {
			if err == sql.ErrNoRows {
				leaguesToImport = append(leaguesToImport, League{
					Id: league.Id,
				})

			} else {
				log.Fatal(err)
			}
		}
		if scannedLeague != 0 {
			fmt.Println("Name:", league.DisplayName, "Id:", league.Id, "✅")
		} else {
			fmt.Println("Name:", league.DisplayName, "Id:", league.Id, "❌")
		}
	}

	fmt.Printf("Importing %d leagues\n", len(leaguesToImport))
	for _, league := range leaguesToImport {
		fmt.Println("Importing league", league.Id)
		leagueWithMatches, err := stratz.GetLeagueMatches(client, ctx, league.Id)
		if err != nil {
			log.Fatal(err)
		}

		err = create_league.CreateLeague(db, &leagueWithMatches)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Imported all leagues")

}

package main

import (
	"dota-pro-db/liquipedia"
	"fmt"
	"log"
)

func main() {
	leagues, err := liquipedia.GetLeagueIds()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(leagues)
}

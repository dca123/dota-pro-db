package main

import (
	"context"
	"database/sql"
	"dota-pro-db/stratz"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Khan/genqlient/graphql"
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
	ctx := context.Background()
	client := graphql.NewClient(URL, &http.Client{
		Transport: &authedTransport{
			wrapped: http.DefaultTransport,
		},
	})
	matchId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Importing data for league: %d\n", matchId)
	league, err := stratz.GetLeagueMatches(client, ctx, matchId)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("League: ", league.DisplayName, "No. of matches: ", len(league.Matches))

	db, err := sql.Open("sqlite3", "../../test.db?_foreign_keys=true")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = createLeague(db, &League{Id: &league.Id, DisplayName: &league.DisplayName})
	if err != nil {
		if isDuplicateError(err) {
			log.Println("Existing league, syncing matches")
		} else {
			log.Fatal(err)
		}
	}

	for _, match := range league.Matches {
		log.Printf("Processing match %d\n", match.Id)
		err := createTeam(db, &Team{Id: &match.RadiantTeam.Id, Name: &match.RadiantTeam.Name, Tag: &match.RadiantTeam.Tag})
		if err != nil {
			if !isDuplicateError(err) {
				log.Fatal(err)
			}
		}
		err = createTeam(db, &Team{Id: &match.DireTeam.Id, Name: &match.DireTeam.Name, Tag: &match.DireTeam.Tag})
		if err != nil {
			if !isDuplicateError(err) {
				log.Fatal(err)
			}
		}

		createSeries(db, &Series{
			Id:              &match.RadiantTeamId,
			Type:            &match.Series.Type,
			TeamOneId:       &match.Series.TeamOneId,
			TeamTwoId:       &match.Series.TeamTwoId,
			WinningTeamId:   &match.Series.WinningTeamId,
			TeamOneWinCount: &match.Series.TeamOneWinCount,
			TeamTwoWinCount: &match.Series.TeamTwoWinCount,
		})
		if err != nil {
			log.Fatalf("Error creating series: %s", err)
		}

	}
}

func isDuplicateError(err error) bool {
	if err, ok := err.(sqlite3.Error); ok {
		if err.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
			return true
		}
		return false
	}
	return false
}

type League = struct {
	Id          *int
	DisplayName *string
}

func createLeague(db *sql.DB, league *League) error {
	stmt, err := db.Prepare("INSERT INTO leagues(id, name) VALUES (?, ?)")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(league.Id, league.DisplayName)
	return err
}

type Team = struct {
	Id   *int
	Name *string
	Tag  *string
}

func createTeam(db *sql.DB, team *Team) error {
	stmt, err := db.Prepare("INSERT INTO teams(id, name, tag) VALUES (?, ?, ?)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(*team.Id, *team.Name, *team.Tag)
	return err
}

type Series = struct {
	Id        *int
	Type      *stratz.SeriesEnum
	TeamOneId *int
	TeamTwoId *int
	//Whats happens if parse in middle of series ?
	TeamOneWinCount *int
	TeamTwoWinCount *int
	WinningTeamId   *int
}

func createSeries(db *sql.DB, series *Series) error {
	stmt, err := db.Prepare("INSERT INTO series(id, type, team_one_win_count, team_two_win_count, winning_team_id, team_one_id, team_two_id) VALUES (?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(*series.Id, *series.Type, *series.TeamOneWinCount, *series.TeamTwoWinCount, *series.WinningTeamId, *series.TeamOneId, *series.TeamTwoId)
	return err
}

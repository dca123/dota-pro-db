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
	"strings"

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

		err = createSeries(db, &Series{
			Id:              &match.Series.Id,
			Type:            &match.Series.Type,
			TeamOneId:       &match.Series.TeamOneId,
			TeamTwoId:       &match.Series.TeamTwoId,
			WinningTeamId:   &match.Series.WinningTeamId,
			TeamOneWinCount: &match.Series.TeamOneWinCount,
			TeamTwoWinCount: &match.Series.TeamTwoWinCount,
		})
		if err != nil {
			if !isDuplicateError(err) {
				log.Fatalf("Error creating series: %s", err)
			} else {
				log.Println("Duplicate series avoided")
			}
		}

		err = createMatch(db, &Match{
			Id:                     &match.Id,
			DidRadiantWin:          &match.DidRadiantWin,
			DurationSeconds:        &match.DurationSeconds,
			StartDateTime:          &match.StartDateTime,
			EndDateTime:            &match.EndDateTime,
			TowerStatusRadiant:     &match.TowerStatusRadiant,
			TowerStatusDire:        &match.TowerStatusDire,
			BarracksStatusRadiant:  &match.BarracksStatusRadiant,
			BarracksStatusDire:     &match.BarracksStatusDire,
			FirstBloodTime:         &match.FirstBloodTime,
			LobbyType:              &match.LobbyType,
			GameMode:               &match.GameMode,
			GameVersionId:          &match.GameVersionId,
			RadiantNetworthLeads:   &match.RadiantNetworthLeads,
			RadiantExperienceLeads: &match.RadiantExperienceLeads,
			AnalysisOutcome:        &match.AnalysisOutcome,
			LeagueId:               &league.Id,
			SeriesId:               &match.SeriesId,
			RadiantTeamId:          &match.RadiantTeamId,
			DireTeamId:             &match.DireTeamId,
		})

		if err != nil {
			if !isDuplicateError(err) {
				log.Fatalf("Error creating match: %s", err)
			} else {
				log.Println("Duplicate match avoided")
			}
		}

		var heroSelections = make([]HeroSelection, len(match.PickBans))
		for i, selection := range match.PickBans {
			heroSelections[i] = HeroSelection{
				IsPick:    &selection.IsPick,
				PickOrder: &selection.Order,
				IsRadiant: &selection.IsRadiant,
				MatchId:   &match.Id,
				HeroId:    &selection.HeroId,
			}
		}

		err = createMatchPickBans(db, &heroSelections)
		if err != nil {
			if !isDuplicateError(err) {
				log.Fatalf("Error creating pick bans: %s", err)
			} else {
				log.Println("Duplicate pick bans avoided")
			}
		}

		var matchPlayers = make([]MatchPlayer, len(match.Players))

		for i, player := range match.Players {
			teamPlayer := TeamPlayer{
				SteamAccountId: &player.SteamAccountId,
				Name:           player.SteamAccount.ProSteamAccount.Name,
			}
			err = createTeamPlayer(db, &teamPlayer)
			if err != nil {
				if !isDuplicateError(err) {
					log.Fatalf("Error creating team player: %s", err)
				} else {
					log.Println("Duplicate team player avoided")
				}
			}

			matchPlayers[i] = MatchPlayer{
				IsRadiant:         &player.IsRadiant,
				IsVictory:         &player.IsVictory,
				Kills:             &player.Kills,
				Deaths:            &player.Deaths,
				Assists:           &player.Assists,
				NumLastHits:       &player.NumLastHits,
				NumDenies:         &player.NumDenies,
				GoldPerMin:        &player.GoldPerMinute,
				Networth:          &player.Networth,
				ExpPerMin:         &player.ExperiencePerMinute,
				Level:             &player.Level,
				GoldSpent:         &player.GoldSpent,
				HeroDamage:        &player.HeroDamage,
				TowerDamage:       &player.TowerDamage,
				HeroHealing:       &player.HeroHealing,
				Lane:              &player.Lane,
				IsRandom:          &player.IsRandom,
				Position:          &player.Position,
				Role:              &player.Role,
				InvisibleSeconds:  &player.InvisibleSeconds,
				DotaPlusHeroLevel: &player.DotaPlusHeroXp,
				MatchId:           &match.Id,
				SteamAccountId:    &player.SteamAccountId,
				HeroId:            &player.HeroId,
				Item0Id:           &player.Item0Id,
				Item1Id:           &player.Item1Id,
				Item2Id:           &player.Item2Id,
				Item3Id:           &player.Item3Id,
				Item4Id:           &player.Item4Id,
				Item5Id:           &player.Item5Id,
				Backpack0Id:       &player.Backpack0Id,
				Backpack1Id:       &player.Backpack1Id,
				Backpack2Id:       &player.Backpack2Id,
				Neutral0Id:        &player.Neutral0Id,
			}
		}

		err = createMatchPlayers(db, &matchPlayers)
		if err != nil {
			if !isDuplicateError(err) {
				log.Fatalf("Error creating match players: %s", err)
			} else {
				log.Println("Duplicate match players avoided")
			}
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
	stmt, err := db.Prepare("INSERT INTO series(id, type, team_one_win_count, team_two_win_count, winning_team_id, team_one_id, team_two_id) VALUES (?, ?, ?, ?, NULLIF(?, 0), ?, ?)")

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(*series.Id, *series.Type, *series.TeamOneWinCount, *series.TeamTwoWinCount, *series.WinningTeamId, *series.TeamOneId, *series.TeamTwoId)
	return err
}

type Match = struct {
	Id                     *int
	DidRadiantWin          *bool
	DurationSeconds        *int
	StartDateTime          *int
	EndDateTime            *int
	TowerStatusRadiant     *int
	TowerStatusDire        *int
	BarracksStatusRadiant  *int
	BarracksStatusDire     *int
	FirstBloodTime         *int
	LobbyType              *stratz.LobbyTypeEnum
	GameMode               *stratz.GameModeEnumType
	GameVersionId          *int
	RadiantNetworthLeads   *[]int
	RadiantExperienceLeads *[]int
	AnalysisOutcome        *stratz.MatchAnalysisOutcomeType
	LeagueId               *int
	SeriesId               *int
	RadiantTeamId          *int
	DireTeamId             *int
}

func boolToInt(b bool) int {
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}

func createMatch(db *sql.DB, match *Match) error {
	stmt, err := db.Prepare("INSERT into matches(id, did_radiant_win, duration_seconds, start_date_time, end_date_time, tower_status_radiant, tower_status_dire, barracks_status_radiant, barracks_status_dire, first_blood_time, lobby_type, game_mode, game_version_id, radiant_networth_leads, radiant_experience_leads, analysis_outcome, league_id, series_id, radiant_team_id, dire_team_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	didRadiantWin := boolToInt(*match.DidRadiantWin)
	radiantNetworthLeads := strings.Join(strings.Fields(fmt.Sprint(*match.RadiantNetworthLeads)), ",")
	radiantExperienceLeads := strings.Join(strings.Fields(fmt.Sprint(*match.RadiantExperienceLeads)), ",")

	_, err = stmt.Exec(
		*match.Id,
		didRadiantWin,
		*match.DurationSeconds,
		*match.StartDateTime,
		*match.EndDateTime,
		*match.TowerStatusRadiant,
		*match.TowerStatusDire,
		*match.BarracksStatusRadiant,
		*match.BarracksStatusDire,
		*match.FirstBloodTime,
		*match.LobbyType,
		*match.GameMode,
		*match.GameVersionId,
		radiantNetworthLeads,
		radiantExperienceLeads,
		*match.AnalysisOutcome,
		*match.LeagueId,
		*match.SeriesId,
		*match.RadiantTeamId,
		*match.DireTeamId,
	)
	return err
}

type HeroSelection = struct {
	IsPick    *bool
	PickOrder *int
	IsRadiant *bool
	MatchId   *int
	HeroId    *int
}

func createMatchPickBans(db *sql.DB, selections *[]HeroSelection) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO match_pick_bans(is_pick, pick_order, is_radiant, match_id, hero_id) VALUES (?, ?, ?, ?, ?)")
	defer stmt.Close()

	for _, selection := range *selections {
		isPick := boolToInt(*selection.IsPick)
		isRadiant := boolToInt(*selection.IsRadiant)
		_, err := stmt.Exec(isPick, *selection.PickOrder, isRadiant, *selection.MatchId, *selection.HeroId)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

type TeamPlayer = struct {
	SteamAccountId *int
	Name           string
}

func createTeamPlayer(db *sql.DB, teamPlayer *TeamPlayer) error {
	stmt, err := db.Prepare("INSERT INTO team_players(steam_account_id, name) VALUES(?,COALESCE(?, ''))")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(teamPlayer.SteamAccountId, teamPlayer.Name)
	return err
}

type MatchPlayer = struct {
	IsRadiant         *bool
	IsVictory         *bool
	Kills             *byte
	Deaths            *byte
	Assists           *byte
	NumLastHits       *int
	NumDenies         *int
	GoldPerMin        *int
	Networth          *int
	ExpPerMin         *int
	Level             *byte
	GoldSpent         *int
	HeroDamage        *int
	TowerDamage       *int
	HeroHealing       *int
	IsRandom          *bool
	Lane              *stratz.MatchLaneType
	Position          *stratz.MatchPlayerPositionType
	Role              *stratz.MatchPlayerRoleType
	InvisibleSeconds  *int
	DotaPlusHeroLevel *int
	MatchId           *int
	SteamAccountId    *int
	HeroId            *int
	Item0Id           *int
	Item1Id           *int
	Item2Id           *int
	Item3Id           *int
	Item4Id           *int
	Item5Id           *int
	Backpack0Id       *int
	Backpack1Id       *int
	Backpack2Id       *int
	Neutral0Id        *int
}

func createMatchPlayers(db *sql.DB, players *[]MatchPlayer) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO match_players(is_radiant, is_victory, kills, deaths, assists, num_last_hits, num_denies, gold_per_min, networth, exp_per_min, level, gold_spent, hero_damage, tower_damage, hero_healing, is_random, lane, position, role, invisible_seconds, dota_plus_hero_level,match_id, steam_account_id, hero_id, item_0_id, item_1_id, item_2_id, item_3_id, item_4_id, item_5_id, backpack_0_id, backpack_1_id, backpack_2_id, neutral_0_id) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,NULLIF(?,0),NULLIF(?,0),NULLIF(?,0),NULLIF(?,0),NULLIF(?,0),NULLIF(?,0),NULLIF(?,0),NULLIF(?,0),NULLIF(?,0),NULLIF(?,0))")
	defer stmt.Close()
	for _, matchPlayer := range *players {
		isRadiant := boolToInt(*matchPlayer.IsRadiant)
		isVictory := boolToInt(*matchPlayer.IsVictory)
		isRandom := boolToInt(*matchPlayer.IsRandom)
		_, err := stmt.Exec(
			isRadiant,
			isVictory,
			*matchPlayer.Kills,
			*matchPlayer.Deaths,
			*matchPlayer.Assists,
			*matchPlayer.NumLastHits,
			*matchPlayer.NumDenies,
			*matchPlayer.GoldPerMin,
			*matchPlayer.Networth,
			*matchPlayer.ExpPerMin,
			*matchPlayer.Level,
			*matchPlayer.GoldSpent,
			*matchPlayer.HeroDamage,
			*matchPlayer.TowerDamage,
			*matchPlayer.HeroHealing,
			isRandom,
			*matchPlayer.Lane,
			*matchPlayer.Position,
			*matchPlayer.Role,
			*matchPlayer.InvisibleSeconds,
			*matchPlayer.DotaPlusHeroLevel,
			*matchPlayer.MatchId,
			*matchPlayer.SteamAccountId,
			*matchPlayer.HeroId,
			*matchPlayer.Item0Id,
			*matchPlayer.Item1Id,
			*matchPlayer.Item2Id,
			*matchPlayer.Item3Id,
			*matchPlayer.Item4Id,
			*matchPlayer.Item5Id,
			*matchPlayer.Backpack0Id,
			*matchPlayer.Backpack1Id,
			*matchPlayer.Backpack2Id,
			*matchPlayer.Neutral0Id,
		)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err

}

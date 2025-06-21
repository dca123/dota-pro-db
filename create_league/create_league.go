package create_league

import (
	"database/sql"
	"dota-pro-db/stratz"
	"fmt"
	"log/slog"
	"strings"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type IngestionError struct {
	Message string
}

func (e *IngestionError) Error() string {
	return fmt.Sprintf("Ingestion error: %s", e.Message)
}

func CreateLeague(db *sql.DB, league *stratz.League) error {
	err := createLeague(db, league)
	if err != nil {
		return err
	}

	for _, match := range league.Matches {
		slog.Info(fmt.Sprintf("Processing match %d\n", match.Id))

		// create radiant team
		err := createTeam(db, &Team{Id: &match.RadiantTeam.Id, Name: &match.RadiantTeam.Name, Tag: &match.RadiantTeam.Tag})
		if err != nil {
			return err
		}

		//create dire team
		err = createTeam(db, &Team{Id: &match.DireTeam.Id, Name: &match.DireTeam.Name, Tag: &match.DireTeam.Tag})
		if err != nil {
			return err
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
			return err
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
			return err

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
			return err
		}

		var matchPlayers = make([]MatchPlayer, len(match.Players))

		for i, player := range match.Players {
			teamPlayer := TeamPlayer{
				SteamAccountId: &player.SteamAccountId,
				Name:           player.SteamAccount.ProSteamAccount.Name,
			}
			err = createTeamPlayer(db, &teamPlayer)
			if err != nil {
				return err
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
			return err
		}

	}
	return nil
}

type League = struct {
	Id          *int
	DisplayName *string
}

func createLeague(db *sql.DB, league *stratz.League) error {
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
	stmt, err := db.Prepare("INSERT INTO teams(id, name, tag) VALUES (?, ?, ?) ON CONFLICT(id) DO NOTHING")

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
	stmt, err := db.Prepare("INSERT INTO series(id, type, team_one_win_count, team_two_win_count, winning_team_id, team_one_id, team_two_id) VALUES (?, ?, ?, ?, NULLIF(?, 0), ?, ?) ON CONFLICT(id) DO NOTHING")

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
	stmt, err := db.Prepare("INSERT INTO team_players(steam_account_id, name) VALUES(?,COALESCE(?, '')) ON CONFLICT (steam_account_id) DO NOTHING")
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

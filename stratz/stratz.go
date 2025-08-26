package stratz

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Khan/genqlient/graphql"
)

type Hero = getHerosConstantsConstantQueryHeroesHeroType
type Item = getItemsConstantsConstantQueryItemsItemType
type League = getLeagueMatchesLeagueLeagueType
type GameVersion = getGameVersionsConstantsConstantQueryGameVersionsGameVersionType

func GetHeroes(client graphql.Client, ctx context.Context) (*getHerosResponse, error) {
	return getHeros(ctx, client)
}

func GetItems(client graphql.Client, ctx context.Context) (*getItemsResponse, error) {
	return getItems(ctx, client)
}

func GetLeagues(client graphql.Client, ctx context.Context, startDateTime int, minPrizePool int) ([]getLeaguesLeaguesLeagueType, error) {
	now := time.Now()
	epoc := now.Unix()
	const TAKE = 100
	resp, err := getLeagues(ctx, client, TAKE, 0, startDateTime, int(epoc))
	if err != nil {
		return []getLeaguesLeaguesLeagueType{}, err
	}
	var leagues []getLeaguesLeaguesLeagueType
	idx := 1
	for len(resp.Leagues) > 0 {
		for _, league := range resp.Leagues {
			if league.PrizePool >= minPrizePool && len(league.Matches) > 9 {
				leagues = append(leagues, league)

			}

		}
		resp, err = getLeagues(ctx, client, TAKE, idx*TAKE, startDateTime, int(epoc))
		if err != nil {
			return []getLeaguesLeaguesLeagueType{}, err
		}
		idx += 1
	}

	return leagues, nil

}

func GetLeagueMatches(client graphql.Client, ctx context.Context, leagueId int) (getLeagueMatchesLeagueLeagueType, error) {
	const TAKE = 100
	resp, err := getLeagueMatches(ctx, client, leagueId, TAKE, 0)
	if err != nil {
		return getLeagueMatchesLeagueLeagueType{}, err
	}

	var matches []getLeagueMatchesLeagueLeagueTypeMatchesMatchType

	idx := 1
	for len(resp.League.Matches) > 0 {
		slog.Debug(fmt.Sprintf("Response with %d matches", len(resp.League.Matches)))
		matches = append(matches, resp.League.Matches...)
		resp, err = getLeagueMatches(ctx, client, leagueId, TAKE, TAKE*idx)
		idx += 1
		if err != nil {
			return getLeagueMatchesLeagueLeagueType{}, err
		}
	}

	var leagueWithAllMatches = getLeagueMatchesLeagueLeagueType{
		Id:          resp.League.Id,
		DisplayName: resp.League.DisplayName,
		Matches:     matches,
	}

	return leagueWithAllMatches, nil
}

func GetGameVersions(client graphql.Client, ctx context.Context) (*getGameVersionsResponse, error) {
	return getGameVersions(ctx, client)
}

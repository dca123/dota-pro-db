package stratz

import (
	"context"
	"fmt"
	"log/slog"

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

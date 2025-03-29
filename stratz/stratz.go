package stratz

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

type Hero = getHerosConstantsConstantQueryHeroesHeroType
type Item = getItemsConstantsConstantQueryItemsItemType

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

	fmt.Println(len(resp.League.DisplayName))

	idx := 1
	for len(resp.League.Matches) > 0 {
		fmt.Println(len(resp.League.Matches))
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

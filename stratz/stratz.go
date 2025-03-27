package stratz

import (
	"context"
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

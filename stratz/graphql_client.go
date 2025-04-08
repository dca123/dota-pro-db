package stratz

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

func GetClient() graphql.Client {
	if API_KEY == "" {
		log.Fatalln("STRATZ_API_KEY is not set")
	}
	client := graphql.NewClient(URL, &http.Client{
		Transport: &authedTransport{
			wrapped: http.DefaultTransport,
		},
	})
	return client
}

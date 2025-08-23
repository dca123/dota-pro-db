package liquipedia

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// League represents a Dota 2 league with its ID and name
type League struct {
	ID   int
	Name string
}

// GetLeagueIds fetches all leagues and their corresponding STRATZ IDs
func GetLeagueIds() ([]League, error) {
	leagues, err := getLeagues()
	if err != nil {
		return nil, fmt.Errorf("failed to get leagues: %w", err)
	}

	var results []League
	for _, league := range leagues {
		id, err := extractStratzID(league.Name)
		if err != nil {
			log.Printf("Warning: failed to get STRATZ ID for %s: %v", league.Name, err)
			continue
		}
		results = append(results, League{
			ID:   id,
			Name: league.Name,
		})
	}

	return results, nil
}

// leagueInfo holds temporary league data
type leagueInfo struct {
	Name string
	URL  string
}

// getLeagues fetches all leagues using cached tier1.txt in dev mode
func getLeagues() ([]leagueInfo, error) {
	// Try to use cached file first (dev mode)
	file, err := os.Open("liquipedia/tier1.txt")
	if err == nil {
		defer file.Close()
		return parseLeaguesFromFile(file)
	}

	// Fallback to live API call
	return getLeaguesFromAPI()
}

// parseLeaguesFromFile parses leagues from cached JSON file
func parseLeaguesFromFile(file *os.File) ([]leagueInfo, error) {
	var result struct {
		Parse struct {
			Text map[string]string `json:"text"`
		} `json:"parse"`
	}

	if err := json.NewDecoder(file).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse cached file: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(result.Parse.Text["*"]))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var leagues []leagueInfo
	doc.Find("div.gridRow").Each(func(i int, row *goquery.Selection) {
		cell := row.Find("div.gridCell.Tournament.Header")
		link := cell.Find("a")
		name := strings.TrimSpace(link.Text())
		href, exists := link.Attr("href")
		if exists && name != "" {
			leagues = append(leagues, leagueInfo{
				Name: name,
				URL:  "https://liquipedia.net" + href,
			})
		}
	})

	return leagues, nil
}

// getLeaguesFromAPI fetches leagues from live API
func getLeaguesFromAPI() ([]leagueInfo, error) {
	endpoint := "https://liquipedia.net/dota2/api.php?action=parse&prop=text&page=Tier_1_Tournaments&format=json"

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "MyBot/1.0 (my@email.com)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Parse struct {
			Text map[string]string `json:"text"`
		} `json:"parse"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(result.Parse.Text["*"]))
	if err != nil {
		return nil, err
	}

	var leagues []leagueInfo
	doc.Find("div.gridRow").Each(func(i int, row *goquery.Selection) {
		cell := row.Find("div.gridCell.Tournament.Header")
		link := cell.Find("a")
		name := strings.TrimSpace(link.Text())
		href, exists := link.Attr("href")
		if exists && name != "" {
			leagues = append(leagues, leagueInfo{
				Name: name,
				URL:  "https://liquipedia.net" + href,
			})
		}
	})

	return leagues, nil
}

// extractStratzID extracts the STRATZ league ID from a Liquipedia tournament page
func extractStratzID(leagueName string) (int, error) {
	endpoint := "https://liquipedia.net/dota2/api.php"

	params := url.Values{}
	params.Set("action", "parse")
	params.Set("prop", "externallinks")
	params.Set("page", leagueName)
	params.Set("format", "json")

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return 0, err
	}
	req.Header.Add("User-Agent", "MyBot/1.0 (my@email.com)")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Parse struct {
			ExternalLinks []string `json:"externallinks"`
		} `json:"parse"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	prefix := "https://stratz.com/leagues/"
	for _, url := range result.Parse.ExternalLinks {
		if strings.HasPrefix(url, prefix) {
			idStr := strings.TrimPrefix(url, prefix)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return 0, fmt.Errorf("invalid ID format: %s", idStr)
			}
			return id, nil
		}
	}

	return 0, fmt.Errorf("no STRATZ ID found for league: %s", leagueName)
}

const rateLimit = time.Second * 30

type Client interface {
	Call(*Payload)
}

type Payload struct {
}

func RateLimitCall(client Client, payloads []*Payload) {
	throttle := time.Tick(rateLimit)
	for _, payload := range payloads {
		<-throttle
		go client.Call(payload)
	}
}

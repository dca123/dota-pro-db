package liquipedia

import (
	"database/sql"
	"dota-pro-db/database"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// League represents a Dota 2 league with its ID and name
type League struct {
	ID   int
	Name string
}

const rateLimit = time.Second * 30

var rateLimiter = time.NewTicker(rateLimit)

var client = &http.Client{Timeout: rateLimit}

func GetLeagueIds() ([]League, error) {
	leagues, err := getLeagues()
	if err != nil {
		return nil, fmt.Errorf("failed to get leagues: %w", err)
	}

	var results []League
	var mu sync.Mutex
	var wg sync.WaitGroup

	rateLimitChan := make(chan struct{}, 1)
	go func() {
		for range rateLimiter.C {
			select {
			case rateLimitChan <- struct{}{}:
			default:
			}
		}
	}()

	for _, league := range leagues {
		wg.Add(1)
		go func(l leagueInfo) {
			defer wg.Done()
			hasLeague, err := doesDBHaveLeague(l.Name)
			if err != nil {
				log.Printf("Warning: failed to check DB for %s: %v", l.Name, err)
				return
			}
			if hasLeague {
				log.Printf("⏭️ League %s already exists in DB, skipping", l.Name)
				return
			}
			log.Printf("⏳ Waiting on %s", league.Name)
			<-rateLimitChan
			log.Printf("⚙️ Processing on %s", league.Name)
			id, err := extractStratzID(league.Name)
			if err != nil {
				log.Printf("Warning: failed to get STRATZ ID for %s: %v", league.Name, err)
				return
			}
			mu.Lock()
			results = append(results, League{
				ID:   id,
				Name: league.Name,
			})
			mu.Unlock()
		}(league)
	}
	wg.Wait()

	return results, nil
}

func doesDBHaveLeague(leagueName string) (bool, error) {
	db := database.GetDb()
	var scannedLeague int
	err := db.QueryRow("SELECT id FROM leagues where liquipedia_page_name = ?", leagueName).Scan(&scannedLeague)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if scannedLeague != 0 {
		return true, nil
	}
	return false, nil
}

// leagueInfo holds temporary league data
type leagueInfo struct {
	Name string
	URL  string
}

func getLeagues() ([]leagueInfo, error) {
	response, err := getLeaguesFromAPI()
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response))
	if err != nil {
		return nil, err
	}

	var leagues []leagueInfo
	doc.Find("div.gridRow").Each(func(i int, row *goquery.Selection) {
		dateCell := row.Find("div.gridCell.Date.Header")
		dateParts := strings.Split(dateCell.Text(), ",")
		year, err := strconv.Atoi(strings.TrimSpace(dateParts[len(dateParts)-1]))
		if err != nil {
			log.Fatalln(err)
		}

		winnerCell := row.Find("div.gridCell.Placement.FirstPlace")
		hasWinner := strings.TrimSpace(winnerCell.Text()) != "TBD"

		if year >= 2025 && hasWinner {
			cell := row.Find("div.gridCell.Tournament.Header")
			link := cell.Find("a").Last()
			name := strings.TrimSpace(link.Text())
			href, exists := link.Attr("href")
			if exists && name != "" {
				leagues = append(leagues, leagueInfo{
					Name: strings.TrimPrefix(href, "/dota2/"),
					URL:  "https://liquipedia.net" + href,
				})
			}

		}

	})
	return leagues, nil
}

func getLeaguesFromFile() (string, error) {
	var result struct {
		Parse struct {
			Text map[string]string `json:"text"`
		} `json:"parse"`
	}
	file, err := os.Open("liquipedia/tier1.txt")
	if err != nil {
		return "", err
	}

	if err := json.NewDecoder(file).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse cached file: %w", err)
	}
	return result.Parse.Text["*"], nil
}

func getLeaguesFromAPI() (string, error) {
	result, err := getLeaguesFromFile()
	if err != nil {
		endpoint := "https://liquipedia.net/dota2/api.php?action=parse&prop=text&page=Tier_1_Tournaments&format=json"
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			return "", err
		}
		req.Header.Add("User-Agent", "MyBot/1.0 (my@email.com)")

		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		var result struct {
			Parse struct {
				Text map[string]string `json:"text"`
			} `json:"parse"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return "", err
		}
		return result.Parse.Text["*"], nil
	}
	return result, nil
}

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

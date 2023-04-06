package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

var total int = 0

// tweet the data out
func tweet(s string) { //} ([]byte, error) {
	// config := oauth1.NewConfig(os.Getenv("TWITTER_API_KEY"), os.Getenv("TWITTER_API_SECRET"))
	// token := oauth1.NewToken(os.Getenv("TWITTER_CLIENT_KEY"), os.Getenv("TWITTER_CLIENT_SECRET"))

	// client := config.Client(oauth1.NoContext, token)

	// fmt.Println(s)
	// twitterUrl := "https://api.twitter.com/2/tweets"

	// payload := strings.NewReader(fmt.Sprintf(`{ "text": "%v" }`, s))

	// fmt.Println(payload)

	// res, err := client.Post(twitterUrl, "application/json", payload)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	// defer res.Body.Close()

	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }

	// fmt.Println(string(body))
	// return body, nil
	fmt.Printf("\n\nTWEETING:\n%v", s)
}

// call the mlb api (or well... the endpoint)
func callAPI(url string) ([]byte, error) {
	client := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := client.Do(req)
	if getErr != nil {
		return make([]byte, 0), getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	return ioutil.ReadAll(res.Body)
}

// parse a scoreboard for gamePks
func parseScoreboardData(gameDay time.Time) (ScheduleData, error) {
	data := ScheduleData{}

	body, err := callAPI(fmt.Sprintf("https://bdfed.stitch.mlbinfra.com/bdfed/transform-mlb-scoreboard?stitch_env=prod&sortTemplate=4&sportId=1&startDate=%v&endDate=%v&gameType=E&&gameType=S&&gameType=R&&gameType=F&&gameType=D&&gameType=L&&gameType=W&&gameType=A&language=en&leagueId=104&&leagueId=103&contextTeamId=", gameDay.Format("2006-01-02"), gameDay.Format("2006-01-02")))

	if err != nil {
		return data, err
	}

	jsonErr := json.Unmarshal(body, &data)

	if jsonErr != nil {
		return data, jsonErr
	}

	return data, nil
}

func main() {
	godotenv.Load()

	gameDay := time.Now()
	// gameDay = gameDay.AddDate(0, 0, -1)
	data, err := parseScoreboardData(gameDay)

	games := []GameData{}
	gamesToCome := []GameData{}
	if err != nil {
		fmt.Printf("Error getting games for day %v", gameDay.Format("Mon Jan 02 2006"))
		log.Fatal(err)
	}

	for _, g := range data.Dates[0].Games {
		game, err := getGameData(g)

		// need to check here if game has started or not
		if err != nil {
			// handle later
		}

		if game.DelayInSeconds > 0 {
			gamesToCome = append(gamesToCome, *game)
		} else {
			games = append(games, *game)
		}
	}

	fmt.Printf("started: %d. not yet: %d", len(games), len(gamesToCome))

	for _, c := range games {
		violations := getGameViolations(c)

		if len(violations) > 0 {
			c.ViolationsTotal += len(violations)
			total += len(violations)
			for _, v := range violations {
				tweet(fmt.Sprintf("%v\n%d in game, %d total in today's games.", v, c.ViolationsTotal, total))
			}
		}
	}

	fmt.Println("Done for the day!!")
}

func getGameData(game ScheduleGame) (*GameData, error) {
	data := GameData{}

	gameId := game.GamePk
	body, _ := callAPI(fmt.Sprintf("https://statsapi.mlb.com/api/v1.1/game/%d/feed/live?language=en", gameId))

	// TODO : fuck with errors later
	// if err != nil {
	// 	return data, err
	// }

	jsonErr := json.Unmarshal(body, &data)

	if jsonErr != nil {
		// return data, jsonErr
		// TODO : fuck with errors later
	}

	data.ViolationsTotal = 0
	data.LastPlayIdx = 0
	data.LastPlayEventIdx = 0

	if time.Now().Before(game.GameDate) {
		data.DelayInSeconds = int(math.Round(game.GameDate.Sub(time.Now()).Seconds()))
	} else {
		data.DelayInSeconds = 0
	}

	return &data, nil
}

func buildNotification(vt string, p Play, e PlayEvent, homeTeam string, awayTeam string) string {
	var a string
	if vt == "B" {
		a = "Batter"
	} else {
		a = "Pitcher"
	}
	t := fmt.Sprintf("Pitch Clock Violation on %v", a)

	o := "Out"

	if p.Count.Outs > 1 {
		o = o + "s"
	}

	return fmt.Sprintf("%v\n%v\n%v %v - %v @ %v\nCount Now %v - %v, %v %v", t, p.Matchup.Batter.FullName, p.About.HalfInning, p.About.Inning, awayTeam, homeTeam, e.Count.Balls, e.Count.Strikes, e.Count.Outs, o)
}

func getGameViolations(game GameData) []string {
	lastPlayIdx := game.LastPlayIdx
	lastPlayEventIdx := game.LastPlayEventIdx

	notifications := []string{}

	for pIdx, v := range game.LiveData.Plays.AllPlays {
		if lastPlayIdx > 0 && lastPlayIdx < pIdx {
			continue
		}

		for eIdx, p := range v.PlayEvents {
			if pIdx == lastPlayIdx && eIdx < lastPlayEventIdx {
				// in case its the same plate appearance we don't want to accidentally send multiple tweets of the same thing
				continue
			}

			switch p.Details.Call.Code {
			case "AC":
				n := buildNotification("B", v, p, game.GameData.Teams.Home.Abbreviation, game.GameData.Teams.Away.Abbreviation)
				notifications = append(notifications, n)
				break

			case "VP":
				n := buildNotification("P", v, p, game.GameData.Teams.Home.Abbreviation, game.GameData.Teams.Away.Abbreviation)
				notifications = append(notifications, n)
				break

			default:
				break
			}
		}
	}

	return notifications
}

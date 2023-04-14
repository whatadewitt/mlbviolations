package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
	"github.com/whatadewitt/mlbviolations/internal"
)

var total int = 0

// tweet the data out
func tweet(v internal.Tweet, reply_id string) ([]byte, error) {
	if reply_id == "" {
		fmt.Printf("\n\nTWEETING:\n%s\n", v.MainTweet)
	} else {
		fmt.Printf("\n\nTWEETING:\n%s\n", v.ReplyTweet)
	}

	config := oauth1.NewConfig(os.Getenv("TWITTER_API_KEY"), os.Getenv("TWITTER_API_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_ACCESS_TOKEN"), os.Getenv("TWITTER_ACCESS_TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	twitterUrl := "https://api.twitter.com/2/tweets"

	var payload *strings.Reader
	if reply_id == "" {
		fmt.Println("SENDING TWEET")
		payload = strings.NewReader(fmt.Sprintf(`{ "text": "%s" }`, v.MainTweet))
	} else {
		fmt.Println("SENDING RESPONSE")
		payload = strings.NewReader(fmt.Sprintf(`{ "text": "%s", "reply": { "in_reply_to_tweet_id": "%s" } }`, v.ReplyTweet, reply_id))
	}

	res, err := httpClient.Post(twitterUrl, "application/json", payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println("payload %v", payload)

	randy := rand.Float32()
	if v.ReplyTweet != "" && randy <= .25 && reply_id == "" {
		// send a reply with additional data 25% of the time
		// but not to the replies...
		fmt.Println("I'm sending a response!")

		var tweetResp internal.TweetResponse
		err = json.Unmarshal(body, &tweetResp)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		return tweet(v, tweetResp.Data.Id)
	} else {
		return nil, nil // body, nil
	}
}

// parse a scoreboard for gamePks
func parseScoreboardData(gameDay time.Time) ([]*internal.TrackedGame, error) {
	data := internal.ScoreboardData{}

	body, err := internal.CallAPI(fmt.Sprintf("https://bdfed.stitch.mlbinfra.com/bdfed/transform-mlb-scoreboard?stitch_env=prod&sortTemplate=4&sportId=1&startDate=%v&endDate=%v&gameType=E&&gameType=S&&gameType=R&&gameType=F&&gameType=D&&gameType=L&&gameType=W&&gameType=A&language=en&leagueId=104&&leagueId=103&contextTeamId=", gameDay.Format("2006-01-02"), gameDay.Format("2006-01-02")))

	if err != nil {
		return nil, err
	}

	jsonErr := json.Unmarshal(body, &data)

	if jsonErr != nil {
		return nil, jsonErr
	}

	if jsonErr != nil {
		return nil, jsonErr
	}

	return data.Dates[0].Games, nil
}

func main() {
	godotenv.Load()

	db, err := sql.Open("mysql",
		fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_TABLE")))

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	gameDay := time.Now()
	gameDay = gameDay.AddDate(0, 0, -3)
	trackedGames, err := parseScoreboardData(gameDay)
	if err != nil {
		fmt.Printf("Error getting games for day %v", gameDay.Format("Mon Jan 02 2006"))
		log.Fatal(err)
	}

	games := make([]*internal.TrackedGame, 0)
	upcomingGames := make([]*internal.TrackedGame, 0)

	for _, game := range trackedGames {
		// if (game.Status.AbstractGameCode == "F" && game.Status.DetailedState == "Postponed") {
		// F here to handle games that have finished
		// (or been delayed today) -- figure it out
		if game.GameDate.Before(time.Now()) {
			// game has started
			games = append(games, game)
		} else {
			// game hasn't started
			upcomingGames = append(upcomingGames, game)
		}
		// }
	}

	gameCount := len(games) + len(upcomingGames)
	fmt.Printf("checking %d games...\n", len(games))

	for {
		newUpcomingGames := make([]*internal.TrackedGame, 0)

		for _, upcoming := range upcomingGames {
			if upcoming.GameDate.Before(time.Now()) {
				games = append(games, upcoming)
			} else {
				newUpcomingGames = append(newUpcomingGames, upcoming)
			}
		}

		upcomingGames = newUpcomingGames

		newActiveGames := make([]*internal.TrackedGame, 0)

		for _, game := range games {
			game.Refresh()
			violations := game.GetViolations(db)

			if len(violations) > 0 {
				game.ViolationsTotal += int32(len(violations))
				for _, v := range violations {
					total++
					tweet(v, "")
				}
			}

			if game.Game.GameData.Status.AbstractGameCode != "F" {
				newActiveGames = append(newActiveGames, game)
			} else {
				gameCount--
			}
			// }
		}

		games = newActiveGames

		fmt.Printf("left to count... %d", gameCount)
		if gameCount == 0 {
			seasonTotal := getViolationCount(db)
			var signOffTweet internal.Tweet
			signOffTweet.MainTweet = fmt.Sprintf(`That's all for today folks!\nTotal Violations Today: %d\nSeason Total:`, total, seasonTotal)
			tweet(signOffTweet, "")
			os.Exit(3)
		}

		fmt.Printf("checking %d games...\n", len(games))

		if len(games) == 0 && len(upcomingGames) > 0 {
			// wait 10 before checking again
			delayInSeconds := int(math.Round(upcomingGames[0].GameDate.Sub(time.Now()).Seconds()))
			fmt.Printf("Next game will start in %v seconds", delayInSeconds)
			time.Sleep(time.Second * time.Duration(delayInSeconds))
		} else {
			time.Sleep(time.Second * 50)
		}
	}
}

func getViolationCount(db *sql.DB) int {
	var seasonTotal int
	err := db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM player_stats")).Scan(&seasonTotal)

	if err != nil {
		log.Fatal(err)
	}

	return seasonTotal
}

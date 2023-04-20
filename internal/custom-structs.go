package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type ScoreboardData struct {
	Dates []ScoreboardDate `json:"dates"`
}

type ScoreboardDate struct {
	Games []*TrackedGame `json:"games"`
}

type TrackedGame struct {
	GamePk           int32      `json:"gamePk"`
	GameDate         time.Time  `json:"gameDate"`
	Game             *GameData  `json:"game_data,omitempty"`
	Status           GameStatus `json:"status,omitempty"`
	ViolationsTotal  int32      `json:"violationsTotal,omitempty"`
	LastPlayIdx      int32      `json:"lastPlayIdx,omitempty"`
	LastPlayEventIdx int32      `json:"lastPlayEventIdx,omitempty"`
}

type Tweet struct {
	MainTweet  string
	ReplyTweet string
}

type TweetResponse struct {
	Data TweetResponseData `json:"data,omitempty"`
}

type TweetResponseData struct {
	Id                     string   `json:"id,omitempty"`
	Edit_history_tweet_ids []string `json:"edit_history_tweet_ids,omitempty"`
	Text                   string   `json:"string,omitempty"`
}

func (t *TrackedGame) Refresh() {
	fmt.Printf("refreshing %d\n", t.GamePk)

	gameId := t.GamePk
	body, _ := CallAPI(fmt.Sprintf("https://statsapi.mlb.com/api/v1.1/game/%d/feed/live?language=en", gameId))
	// TODO : fuck with errors later
	// if err != nil {
	// 	return data, err
	// }

	jsonErr := json.Unmarshal(body, &t.Game)
	if jsonErr != nil {
		// return data, jsonErr
		// TODO : fuck with errors later
	}
}

func (t *TrackedGame) GetViolations(db *sql.DB) []Tweet {
	// Request the data, and return just the violations.
	lastPlayIdx := t.LastPlayIdx
	lastPlayEventIdx := t.LastPlayEventIdx

	notifications := []Tweet{}

	// fmt.Printf("start @ %d %d\n\n", lastPlayIdx, lastPlayEventIdx)

	for pIdx, v := range t.Game.LiveData.Plays.AllPlays {
		if lastPlayIdx > 0 && int(lastPlayIdx) > pIdx {
			continue
		}
		// fmt.Printf("i: %d\n", pIdx)

		for eIdx, p := range v.PlayEvents {
			if pIdx == int(lastPlayIdx) && int(lastPlayEventIdx) >= eIdx {
				// in case its the same plate appearance we don't want to accidentally
				// send multiple tweets of the same thing
				continue
			}
			// fmt.Printf("ie: %d-%d\n", pIdx, eIdx)

			switch p.Details.Call.Code {
			case "AC":
				n := buildNotification(db, "B", v, p, t.Game.GameData.Teams, t.Game.GameData.Datetime.OfficialDate)
				notifications = append(notifications, n)
			case "VP":
				n := buildNotification(db, "P", v, p, t.Game.GameData.Teams, t.Game.GameData.Datetime.OfficialDate)
				notifications = append(notifications, n)
			}

			t.LastPlayEventIdx = int32(eIdx)
		}

		t.LastPlayIdx = int32(pIdx)
	}
	// fmt.Printf("Setting Last Played to: %d-%d", t.LastPlayIdx, t.LastPlayEventIdx)

	// fmt.Printf("L: %d\n", t.LastPlayIdx)
	return notifications
}

func buildNotification(db *sql.DB, vt string, p Play, e PlayEvent, teams GameTeams, date string) Tweet {
	var a string
	var t string
	var tweet Tweet
	var teamAbbr string

	homeTeam := teams.Home.Abbreviation
	awayTeam := teams.Away.Abbreviation
	homeAbbr := teams.Home.TeamName
	awayAbbr := teams.Away.TeamName

	var player Player
	if vt == "B" {
		a = "Batter"
		player = p.Matchup.Batter
		if p.About.HalfInning == "top" {
			t = awayTeam
			teamAbbr = awayAbbr
		} else {
			t = homeTeam
			teamAbbr = homeAbbr
		}
	} else {
		a = "Pitcher"
		player = p.Matchup.Pitcher
		if p.About.HalfInning == "top" {
			t = homeTeam
			teamAbbr = homeAbbr
		} else {
			t = awayTeam
			teamAbbr = awayAbbr
		}
	}
	tt, ttt, pt, ptt := saveViolation(db, a, player, t, date)
	thead := fmt.Sprintf("Pitch Clock Violation on %s", a)

	fmt.Printf("%s: %d, %d, %d, %d\n", player.FullName, tt, ttt, pt, ptt)

	o := "Out"

	if p.Count.Outs > 1 {
		o = o + "s"
	}

	tweet.MainTweet = fmt.Sprintf(`%v\n%v (%v)\n%v %v - %v @ %v\nCount Now %v - %v, %v %v`, thead, player.FullName, t, p.About.HalfInning, p.About.Inning, awayTeam, homeTeam, e.Count.Balls, e.Count.Strikes, e.Count.Outs, o)

	// move to a function
	pluralization := "'"
	if player.FullName[len(player.FullName)-1:] != "s" {
		pluralization = fmt.Sprintf(`%ss`, pluralization)
	}

	suffix := getOrdinalSuffix(pt)

	var dailyNote string
	if ptt > 1 {
		dailyNote = fmt.Sprintf(" (%d%s today)", ptt, getOrdinalSuffix(ptt))
		fmt.Printf("DAILY NOTE!!\n%s\n", dailyNote)
	}

	tweet.ReplyTweet = fmt.Sprintf(`That's %s%s %d%s violation this season%s. As a team the %s have %d violations this year (%d %s violations)`, player.FullName, pluralization, pt, suffix, dailyNote, teamAbbr, tt, ttt, strings.ToLower(a))

	return tweet
}

func getOrdinalSuffix(i int) string {
	j := i % 10
	k := i % 100
	if j == 1 && k != 11 {
		return "st"
	}
	if j == 2 && k != 12 {
		return "nd"
	}
	if j == 3 && k != 13 {
		return "rd"
	}
	return "th"
}

func saveViolation(db *sql.DB, code string, p Player, team string, date string) (int, int, int, int) {
	table := os.Getenv("DB_TABLE")
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s(player_id, player_name, team_name, code, date) VALUES(?, ?, ?, ?, ?)", table))
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(p.ID, p.FullName, team, code, date)
	if err != nil {
		log.Fatal(err)
	}

	// grab the teams total
	var teamTotal int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE team_name = '%s'", table, team)).Scan(&teamTotal)
	if err != nil {
		log.Fatal(err)
	}

	// grab the team total for the given code
	var teamTotalType int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE team_name = '%s' AND code = '%s'", table, team, code)).Scan(&teamTotalType)
	if err != nil {
		log.Fatal(err)
	}

	// grab the player's total
	var playerTotal int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE player_id = %d", table, p.ID)).Scan(&playerTotal)
	if err != nil {
		log.Fatal(err)
	}

	// grab the player's total today
	var playerTotalToday int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE player_id = %d AND date = '%s'", table, p.ID, date)).Scan(&playerTotalToday)
	if err != nil {
		log.Fatal(err)
	}

	return teamTotal, teamTotalType, playerTotal, playerTotalToday
}

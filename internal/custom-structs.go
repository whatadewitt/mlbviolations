package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	ReplyTweet *string
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

func (t *TrackedGame) GetViolations(db *sql.DB) [][2]string {
	// Request the data, and return just the violations.
	lastPlayIdx := t.LastPlayIdx
	lastPlayEventIdx := t.LastPlayEventIdx

	notifications := [][2]string{}

	fmt.Printf("start @ %d %d\n\n", lastPlayIdx, lastPlayEventIdx)

	for pIdx, v := range t.Game.LiveData.Plays.AllPlays {
		if lastPlayIdx > 0 && int(lastPlayIdx) > pIdx {
			continue
		}
		fmt.Printf("i: %d\n", pIdx)

		for eIdx, p := range v.PlayEvents {
			if pIdx == int(lastPlayIdx) && int(lastPlayEventIdx) >= eIdx {
				// in case its the same plate appearance we don't want to accidentally
				// send multiple tweets of the same thing
				continue
			}
			fmt.Printf("ie: %d-%d\n", pIdx, eIdx)

			switch p.Details.Call.Code {
			case "AC":
				n := buildNotification(db, "B", v, p, t.Game.GameData.Teams)
				notifications = append(notifications, n)
			case "VP":
				n := buildNotification(db, "P", v, p, t.Game.GameData.Teams)
				notifications = append(notifications, n)
			}

			t.LastPlayEventIdx = int32(eIdx)
		}

		t.LastPlayIdx = int32(pIdx)

	}
	fmt.Printf("Setting Last Played to: %d-%d", t.LastPlayIdx, t.LastPlayEventIdx)

	// fmt.Printf("L: %d\n", t.LastPlayIdx)
	return notifications
}

func buildNotification(db *sql.DB, vt string, p Play, e PlayEvent, teams GameTeams) Tweet {
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
	tt, ttt, pt := saveViolation(db, a, player, t)
	thead := fmt.Sprintf("Pitch Clock Violation on %s", a)

	o := "Out"

	if p.Count.Outs > 1 {
		o = o + "s"
	}

	tweet.MainTweet = fmt.Sprintf(`%v\n%v (%v)\n%v %v - %v @ %v\nCount Now %v - %v, %v %v`, thead, player.FullName, t, p.About.HalfInning, p.About.Inning, awayTeam, homeTeam, e.Count.Balls, e.Count.Strikes, e.Count.Outs, o)

	pluralization := "'"
	if player.FullName[len(player.FullName)-1:] != "s" {
		pluralization = fmt.Sprintf(`%ss`, pluralization)
	}

	if p.About.HalfInning == "top" {
		teamAbbr = teams.Home.ClubName
	} else {
		teamAbbr = homeTeam
	}

	tweet.ReplyTweet = fmt.Sprintf(`That's %s%s %d%s violation this season. As a team the %s have %d violations this year (%d %s violations)`, player.FullName, pluralization, pt, "th", teamAbbr, tt, ttt, strings.ToLower(a))

	return tweet
}

func saveViolation(db *sql.DB, code string, p Player, team string) (int, int, int) {
	stmt, err := db.Prepare("INSERT INTO player_stats(player_id, player_name, team_name, code, date) VALUES(?, ?, ?, ?, NOW())")
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(p.ID, p.FullName, team, code)
	if err != nil {
		log.Fatal(err)
	}

	// grab the teams total
	var teamTotal int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM player_stats WHERE team_name = '%s'", team)).Scan(&teamTotal)
	if err != nil {
		log.Fatal(err)
	}

	// grab the team total for the given code
	var teamTotalType int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM player_stats WHERE team_name = '%s' AND code = '%s'", team, code)).Scan(&teamTotalType)
	if err != nil {
		log.Fatal(err)
	}

	// grab the player's total
	var playerTotal int
	err = db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM player_stats WHERE player_id = '%d'", p.ID)).Scan(&playerTotal)
	if err != nil {
		log.Fatal(err)
	}

	return teamTotal, teamTotalType, playerTotal
}

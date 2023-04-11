package internal

import (
	"encoding/json"
	"fmt"
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

func (t *TrackedGame) GetViolations() []string {
	// Request the data, and return just the violations.
	lastPlayIdx := t.LastPlayIdx
	lastPlayEventIdx := t.LastPlayEventIdx

	notifications := []string{}

	// fmt.Printf("start @ %d %d\n\n", lastPlayIdx, lastPlayEventIdx)

	for pIdx, v := range t.Game.LiveData.Plays.AllPlays {
		if lastPlayIdx > 0 && int(lastPlayIdx) > pIdx {
			continue
		}
		fmt.Printf("i: %d\n", pIdx)

		for eIdx, p := range v.PlayEvents {
			if pIdx == int(lastPlayIdx) && int(lastPlayEventIdx) > eIdx {
				// in case its the same plate appearance we don't want to accidentally
				// send multiple tweets of the same thing
				continue
			}
			fmt.Printf("ie: %d-%d\n", pIdx, eIdx)

			switch p.Details.Call.Code {
			case "AC":
				n := buildNotification("B", v, p, t.Game.GameData.Teams.Home.Abbreviation, t.Game.GameData.Teams.Away.Abbreviation)
				notifications = append(notifications, n)
			case "VP":
				n := buildNotification("P", v, p, t.Game.GameData.Teams.Home.Abbreviation, t.Game.GameData.Teams.Away.Abbreviation)
				notifications = append(notifications, n)
			}

			t.LastPlayEventIdx = int32(eIdx)
		}

		t.LastPlayIdx = int32(pIdx)

	}
	fmt.Printf("Setting Last Played to: %v-%v", t.LastPlayIdx, t.LastPlayEventIdx)

	// fmt.Printf("L: %d\n", t.LastPlayIdx)
	return notifications
}

func buildNotification(vt string, p Play, e PlayEvent, homeTeam string, awayTeam string) string {
	var a string
	var t string
	var player string
	if vt == "B" {
		a = "Batter"
		player = p.Matchup.Batter.FullName
		if p.About.HalfInning == "top" {
			t = awayTeam
		} else {
			t = homeTeam
		}
	} else {
		a = "Pitcher"
		player = p.Matchup.Pitcher.FullName
		if p.About.HalfInning == "top" {
			t = homeTeam
		} else {
			t = awayTeam
		}
	}
	tweet := fmt.Sprintf("Pitch Clock Violation on %v", a)

	o := "Out"

	if p.Count.Outs > 1 {
		o = o + "s"
	}

	return fmt.Sprintf(`%v\n%v (%v)\n%v %v - %v @ %v\nCount Now %v - %v, %v %v`, tweet, player, t, p.About.HalfInning, p.About.Inning, awayTeam, homeTeam, e.Count.Balls, e.Count.Strikes, e.Count.Outs, o)
}

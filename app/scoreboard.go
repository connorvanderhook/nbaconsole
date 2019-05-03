package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/cnnrrss/nbaconsole/api"
	"github.com/cnnrrss/nbaconsole/common/pad"
)

var (
	scoreBoard        sync.Mutex
	scoreBoardHeaders = []string{"Home", "Away", "Score", "Status"}
)

// getScoreboard locks the current scoreBoard view and requests new data from
// the NBA API. If the data is received, the scoreBoard view is written to stdout
func (nba *NBAConsole) getScoreboard() error {
	nba.g.SetCurrentView(scoreboardLabel)

	params := genericParams(nba.date)
	resp, err := api.GetDataScoreBoard(params)
	if err != nil {
		return fmt.Errorf("Error with request %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading request body %v", err)
	}

	sb := api.DataScoreboard{}
	json.Unmarshal(body, &sb)

	scoreBoard.Lock()
	nba.update(func() {
		curW, _ := nba.g.Size()
		nba.scoreboard.Clear()
		nba.drawScoreboard(nba.scoreboard, sb, curW)
	})
	scoreBoard.Unlock()

	return nil
}

func (nba *NBAConsole) drawScoreboard(output io.Writer, sb api.DataScoreboard, width int) {
	nba.gamesList.Wipe()
	fmt.Fprintln(output, formatScoreBoardHeader(width))

	for _, gm := range sb.Games {
		var blob strings.Builder
		hScore, vScore := gm.Score()
		blob.WriteString(pad.Left(gm.VTeam.TriCode, 5, " "))
		blob.WriteString(pad.Left(gm.HTeam.TriCode, 7, " "))
		blob.WriteString(pad.Left(fmt.Sprintf("%s - %s", vScore, hScore), blob.Len(), " "))
		blob.WriteString(pad.Left(gm.Status(), 14, " "))
		if gm.Playoffs.RoundNum != "" {
			blob.WriteString(pad.Left(gm.Playoffs.SeriesSummaryText, 6, " "))
		}
		fmt.Fprintln(output, blob.String())
		nba.gamesList.games = append(nba.gamesList.games, gm.GameID)
	}
	if len(nba.gamesList.games) == 0 {
		fmt.Fprintf(output, "No hoops today, %s\n", nba.message)
	}

	nba.gamesList.currentIdx = 0 // TODO: implement?

	nba.scoreboard.SetOrigin(0, 0)
	nba.scoreboard.SetCursor(0, 2)

	highlightView(nba.scoreboard)
}

func formatScoreBoardHeader(width int) string {
	var str strings.Builder

	// TODO: this sucks
	for i, h := range scoreBoardHeaders {
		switch {
		case i < len(scoreBoardHeaders)-1:
			str.WriteString(pad.Left(h, 6+(1*i), " "))
		default:
			str.WriteString(pad.Left(h, 17, " "))
		}
	}

	length := str.Len()
	str.WriteString(fmt.Sprintf("\n%s", pad.AddString(length, "-")))

	return str.String()
}

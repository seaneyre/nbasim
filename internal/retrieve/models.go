package retrieve

import "time"

type PlayByPlayResponse struct {
	Meta struct {
		Version int    `json:"version"`
		Code    int    `json:"code"`
		Request string `json:"request"`
		Time    string `json:"time"`
	} `json:"meta"`
	Game struct {
		GameID  string `json:"gameId"`
		Actions []interface{} `json:"actions"`
	} `json:"game"`
}

type GameStart struct {
	ActionNumber    int       `json:"actionNumber"`
	Clock           string    `json:"clock"`
	TimeActual      time.Time `json:"timeActual"`
	Period          int       `json:"period"`
	PeriodType      string    `json:"periodType"`
	ActionType      string    `json:"actionType"`
	SubType         string    `json:"subType"`
	Qualifiers      []any     `json:"qualifiers"`
	PersonID        int       `json:"personId"`
	X               any       `json:"x"`
	Y               any       `json:"y"`
	Possession      int       `json:"possession"`
	ScoreHome       string    `json:"scoreHome"`
	ScoreAway       string    `json:"scoreAway"`
	Edited          time.Time `json:"edited"`
	OrderNumber     int       `json:"orderNumber"`
	XLegacy         any       `json:"xLegacy"`
	YLegacy         any       `json:"yLegacy"`
	IsFieldGoal     int       `json:"isFieldGoal"`
	Side            any       `json:"side"`
	Description     string    `json:"description"`
	PersonIdsFilter []any     `json:"personIdsFilter"`
}

type Jumpball struct {
	ActionNumber             int       `json:"actionNumber"`
	Clock                    string    `json:"clock"`
	TimeActual               time.Time `json:"timeActual"`
	Period                   int       `json:"period"`
	PeriodType               string    `json:"periodType"`
	TeamID                   int       `json:"teamId"`
	TeamTricode              string    `json:"teamTricode"`
	ActionType               string    `json:"actionType"`
	SubType                  string    `json:"subType"`
	Descriptor               string    `json:"descriptor"`
	Qualifiers               []any     `json:"qualifiers"`
	PersonID                 int       `json:"personId"`
	X                        any       `json:"x"`
	Y                        any       `json:"y"`
	Possession               int       `json:"possession"`
	ScoreHome                string    `json:"scoreHome"`
	ScoreAway                string    `json:"scoreAway"`
	Edited                   time.Time `json:"edited"`
	OrderNumber              int       `json:"orderNumber"`
	XLegacy                  any       `json:"xLegacy"`
	YLegacy                  any       `json:"yLegacy"`
	IsFieldGoal              int       `json:"isFieldGoal"`
	JumpBallRecoveredName    string    `json:"jumpBallRecoveredName"`
	JumpBallRecoverdPersonID int       `json:"jumpBallRecoverdPersonId"`
	Side                     any       `json:"side"`
	PlayerName               string    `json:"playerName"`
	PlayerNameI              string    `json:"playerNameI"`
	PersonIdsFilter          []int     `json:"personIdsFilter"`
	JumpBallWonPlayerName    string    `json:"jumpBallWonPlayerName"`
	JumpBallWonPersonID      int       `json:"jumpBallWonPersonId"`
	Description              string    `json:"description"`
	JumpBallLostPlayerName   string    `json:"jumpBallLostPlayerName"`
	JumpBallLostPersonID     int       `json:"jumpBallLostPersonId"`
}
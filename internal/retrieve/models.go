package retrieve

import "time"

type Action struct {
	ActionNumber             int       `json:"actionNumber"`
	Clock                    string    `json:"clock"`
	TimeActual               time.Time `json:"timeActual"`
	Period                   int       `json:"period"`
	PeriodType               string    `json:"periodType"`
	ActionType               string    `json:"actionType"`
	SubType                  string    `json:"subType,omitempty"`
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
	Side                     any       `json:"side"`
	Description              string    `json:"description,omitempty"`
	PersonIdsFilter          []any     `json:"personIdsFilter"`
	TeamID                   int       `json:"teamId,omitempty"`
	TeamTricode              string    `json:"teamTricode,omitempty"`
	Descriptor               string    `json:"descriptor,omitempty"`
	JumpBallRecoveredName    string    `json:"jumpBallRecoveredName,omitempty"`
	JumpBallRecoverdPersonID int       `json:"jumpBallRecoverdPersonId,omitempty"`
	PlayerName               string    `json:"playerName,omitempty"`
	PlayerNameI              string    `json:"playerNameI,omitempty"`
	JumpBallWonPlayerName    string    `json:"jumpBallWonPlayerName,omitempty"`
	JumpBallWonPersonID      int       `json:"jumpBallWonPersonId,omitempty"`
	JumpBallLostPlayerName   string    `json:"jumpBallLostPlayerName,omitempty"`
	JumpBallLostPersonID     int       `json:"jumpBallLostPersonId,omitempty"`
	ShotDistance             float64   `json:"shotDistance,omitempty"`
	ShotResult               string    `json:"shotResult,omitempty"`
	PointsTotal              int       `json:"pointsTotal,omitempty"`
	AssistPlayerNameInitial  string    `json:"assistPlayerNameInitial,omitempty"`
	AssistPersonID           int       `json:"assistPersonId,omitempty"`
	AssistTotal              int       `json:"assistTotal,omitempty"`
	OfficialID               int       `json:"officialId,omitempty"`
	ShotActionNumber         int       `json:"shotActionNumber,omitempty"`
	ReboundTotal             int       `json:"reboundTotal,omitempty"`
	ReboundDefensiveTotal    int       `json:"reboundDefensiveTotal,omitempty"`
	ReboundOffensiveTotal    int       `json:"reboundOffensiveTotal,omitempty"`
	FoulPersonalTotal        int       `json:"foulPersonalTotal,omitempty"`
	FoulTechnicalTotal       int       `json:"foulTechnicalTotal,omitempty"`
	FoulDrawnPlayerName      string    `json:"foulDrawnPlayerName,omitempty"`
	FoulDrawnPersonID        int       `json:"foulDrawnPersonId,omitempty"`
	TurnoverTotal            int       `json:"turnoverTotal,omitempty"`
	StealPlayerName          string    `json:"stealPlayerName,omitempty"`
	StealPersonID            int       `json:"stealPersonId,omitempty"`
	Value                    string    `json:"value,omitempty"`
	BlockPlayerName          string    `json:"blockPlayerName,omitempty"`
	BlockPersonID            int       `json:"blockPersonId,omitempty"`
}

type PlayByPlayResponse struct {
	Meta struct {
		Version int    `json:"version"`
		Code    int    `json:"code"`
		Request string `json:"request"`
		Time    string `json:"time"`
	} `json:"meta"`
	Game struct {
		GameID  string   `json:"gameId"`
		Actions []Action `json:"actions"`
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

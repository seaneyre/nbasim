package simulation

type Event struct {
	GameClockTime int
	ActionType    string
	Action        interface{}
}

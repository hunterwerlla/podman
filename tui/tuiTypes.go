package tui

//view constants

//go:generate stringer -type=TuiScreen
type Screen int

const (
	Subscribed Screen = iota
	Podcast    Screen = iota
	Search     Screen = iota
	Downloaded Screen = iota
)

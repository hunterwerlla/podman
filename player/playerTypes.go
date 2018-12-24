package player

//go:generate stringer -type=PlayerState
type PlayerState int

const (
	// NothingPlaying is the state when the player has nothing cued up/paused/buffered
	NothingPlaying PlayerState = iota
	Resume         PlayerState = iota
	Play           PlayerState = iota
	Pause          PlayerState = iota
	Stop           PlayerState = iota
	FastForward    PlayerState = iota
	Rewind         PlayerState = iota
	ExitPlayer     PlayerState = iota
)

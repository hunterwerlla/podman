package main

import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
)

//go:generate stringer -type=PlayerState
type PlayerState int

const (
	// PlayerNothingPlaying is the state when the player has nothing cued up/paused/buffered
	PlayerNothingPlaying PlayerState = iota
	PlayerResume         PlayerState = iota
	PlayerPlay           PlayerState = iota
	PlayerPause          PlayerState = iota
	PlayerStop           PlayerState = iota
	PlayerFastForward    PlayerState = iota
	PlayerRewind         PlayerState = iota
	PlayerExit           PlayerState = iota
)

var (
	playerPosition int = -1
	startTime      time.Time
	lengthOfFile   int
	playerControl  chan PlayerState
	fileChannel    chan string
	exitChannel    chan bool
	playing        string
	playerState    = PlayerNothingPlaying
)

// StartPlayer starts the global player. The player is global since there is only one of them
// I'm not really a fan of moving it into an object as it should not be reused
func StartPlayer() {
	playerControl = make(chan PlayerState)
	fileChannel = make(chan string)
	exitChannel = make(chan bool)
	go startPlayer()
}

// DisposePlayer sends a signal to the player to destroy itself, and then waits for the player to exit
// This function will deadlock if called twice
func DisposePlayer() {
	playerControl <- PlayerExit
	<-exitChannel
}

func TogglePlayerState() {
	if playerState == PlayerPlay {
		playerControl <- PlayerPause
	} else if playerState == PlayerPause {
		playerControl <- PlayerResume
	}
}

func StopPlayer() {
	playerState = PlayerStop
	playerPosition = -1
	playerState = PlayerNothingPlaying
	playing = ""
	lengthOfFile = 0 //set length
}

func GetLengthOfPlayingFile() int {
	return lengthOfFile
}

func GetPlayerState() PlayerState {
	return playerState
}

func SetPlayerState(state PlayerState) {
	playerControl <- state
}

func SetPlaying(filename string) {
	fileChannel <- filename
}

func GetPlayerPosition() int {
	if playerPosition < 0 {
		return playerPosition
	}
	return playerPosition + int(time.Since(startTime).Seconds())
}

//this runs on its own thread to start/stop and select the media that is playing
func startPlayer() {
	var (
		status                = PlayerNothingPlaying
		stopToExit            = false
		ctrl       *beep.Ctrl = nil
	)
	for stopToExit != true {
		status = PlayerNothingPlaying //reset status

		select {
		case status = <-playerControl:
		case playing = <-fileChannel:
			continue
		}

		switch status {
		case PlayerNothingPlaying:
			panic("invalid state sent to player, this should never happen")
		case PlayerPlay:
			ctrl = playFile()
		case PlayerResume:
			playerState = PlayerPlay
			ctrl.Paused = false
			startTime = time.Now()
		case PlayerPause:
			//save time and file then cleanup
			ctrl.Paused = true
			playerPosition += int(time.Since(startTime).Seconds())
			playerState = PlayerPause
		case PlayerStop:
			ctrl.Paused = true
			ctrl.Streamer = nil
			StopPlayer()
		case PlayerFastForward:
			fastForwardPlayer(ctrl)
		case PlayerRewind:
			rewindPlayer(ctrl)
		case PlayerExit:
			speaker.Clear()
			stopToExit = true
		}
	}
	exitChannel <- true
}

func changePlayerPosition(inputFile beep.StreamSeekCloser, format beep.Format) {
	if playerPosition < 0 {
		playerPosition = 0
	}
	seek := int(float64(playerPosition)*float64(format.SampleRate)*float64(format.NumChannels) + 0.5)
	seek -= seek % int(format.NumChannels)
	_ = inputFile.Seek(seek)
}

func playFile() *beep.Ctrl {
	playerState = PlayerPlay
	inputFile, _ := os.Open(playing)
	decoaded, format, _ := mp3.Decode(inputFile)
	ctrl := &beep.Ctrl{decoaded, false}
	changePlayerPosition(decoaded, format)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(ctrl, beep.Callback(StopPlayer)))
	lengthOfFile = decoaded.Len() / format.NumChannels / int(format.SampleRate)
	startTime = time.Now()
	return ctrl
}

func fastForwardPlayer(ctrl *beep.Ctrl) {
	// TODO rewrite this
}

func rewindPlayer(ctrl *beep.Ctrl) {
	// TODO rewrite this
}

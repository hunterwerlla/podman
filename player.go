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
	playerPosition int = 0
	// TODO get rid of this and just use file size
	startTime     time.Time
	lengthOfFile  int
	playerControl chan PlayerState
	fileChannel   chan string
	exitChannel   chan bool
	playing       string
	playerState   = PlayerNothingPlaying
)

// StartPlayer starts the global player. The player is global since there is only one of them
// I'm not really a fan of moving it into an object as it should not be reused
func StartPlayer(configuration *Configuration) {
	playerControl = make(chan PlayerState)
	fileChannel = make(chan string)
	exitChannel = make(chan bool)
	go startPlayer(configuration)
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
	playing = ""
	lengthOfFile = 0 //set length
}

func GetLengthOfPlayingFile() int {
	return lengthOfFile
}

func GetPlayerState() PlayerState {
	return playerState
}

func sendPlayerMessage(state PlayerState) {
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
func startPlayer(configuration *Configuration) {
	var (
		status                  = PlayerNothingPlaying
		stopToExit              = false
		ctrl       *beep.Ctrl   = nil
		format     *beep.Format = nil
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
			playerState = PlayerNothingPlaying
		case PlayerPlay:
			if playerState != PlayerPause {
				playerPosition = 0
			}
			ctrl, format = playFile()
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
			if playerState == PlayerPlay {
				fastForwardPlayer(ctrl, format, configuration)
			}
		case PlayerRewind:
			if playerState == PlayerPlay {
				rewindPlayer(ctrl, format, configuration)
			}
		case PlayerExit:
			speaker.Clear()
			stopToExit = true
		}
	}
	exitChannel <- true
}

func changePlayerPosition(inputFile beep.StreamSeekCloser, format *beep.Format, position int) {
	if position < 0 {
		position = 0
	}
	frames := format.SampleRate.N(time.Second * time.Duration(position))
	if frames < inputFile.Len() {
		_ = inputFile.Seek(frames)
	} else {
		_ = inputFile.Seek(inputFile.Len() - 1)
	}
}

func playFile() (*beep.Ctrl, *beep.Format) {
	playerState = PlayerPlay
	inputFile, _ := os.Open(playing)
	decodedMp3, format, _ := mp3.Decode(inputFile)
	ctrl := &beep.Ctrl{Streamer: decodedMp3}
	changePlayerPosition(decodedMp3, &format, playerPosition)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(ctrl, beep.Callback(StopPlayer)))
	lengthOfFile = int(format.SampleRate.D(decodedMp3.Len()).Seconds())
	startTime = time.Now()
	return ctrl, &format
}

func fastForwardPlayer(ctrl *beep.Ctrl, format *beep.Format, configuration *Configuration) {
	streamer := ctrl.Streamer.(beep.StreamSeekCloser)
	playerPosition = playerPosition + configuration.FastForwardLength
	changePlayerPosition(streamer, format, playerPosition)
	startTime = time.Now()
}

func rewindPlayer(ctrl *beep.Ctrl, format *beep.Format, configuration *Configuration) {
	streamer := ctrl.Streamer.(beep.StreamSeekCloser)
	playerPosition = playerPosition - configuration.RewindLength
	if playerPosition < 1 {
		playerPosition = 0
	}
	changePlayerPosition(streamer, format, playerPosition)
	startTime = time.Now()
}

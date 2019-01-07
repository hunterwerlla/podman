package main

import (
	"fmt"
	"github.com/krig/go-sox"
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
	lengthOfFile   uint64
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
		playerControl <- PlayerPlay
	}
}

func StopPlayer() {
	//reset position
	playerPosition = -1
	playerState = PlayerNothingPlaying
	playing = ""
	lengthOfFile = 0 //set length
}

func GetLengthOfPlayingFile() uint64 {
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
		chain      *sox.EffectsChain
		inputFile  *sox.Format
		status     = PlayerNothingPlaying
		stopToExit = false
	)
	if !sox.Init() {
		panic("Unable to start the player")
	}
	defer sox.Quit()
	for stopToExit != true {
		status = PlayerNothingPlaying //reset status

		select {
		case status = <-playerControl:
		case playing = <-fileChannel:
			continue
		}

		switch status {
		case PlayerNothingPlaying:
			panic("invalid state when switching status, this should never happen")
		case PlayerPlay:
			if chain != nil || inputFile != nil {
				cleanupSoxData(&chain, &inputFile)
			}
			// TODO yes we are leaking fd's. sox is making pulseaudio panic and we need to get rid of it.
			chain, inputFile, _ = playFile()
		case PlayerPause:
			//save time and file then cleanup
			playerPosition += int(time.Since(startTime).Seconds())
			cleanupSoxData(&chain, &inputFile)
			playerState = PlayerPause
		case PlayerStop:
			StopPlayer()
			cleanupSoxData(&chain, &inputFile)
		case PlayerFastForward:
			fastForwardPlayer()
		case PlayerRewind:
			rewindPlayer()
		case PlayerExit:
			cleanupSoxData(&chain, &inputFile)
			stopToExit = true
		}
	}
	exitChannel <- true
}

func changePlayerPosition(inputFile *sox.Format) {
	if playerPosition < 0 {
		playerPosition = 0
	}
	//formula taken from example 2 of goSoX
	seek := uint64(float64(playerPosition)*float64(inputFile.Signal().Rate())*float64((inputFile.Signal().Channels())) + 0.5)
	seek -= seek % uint64(inputFile.Signal().Channels())
	inputFile.Seek(seek)
}

func playFile() (chain *sox.EffectsChain, inputFile *sox.Format, outputFile *sox.Format) {
	playerState = PlayerPlay
	inputFile = sox.OpenRead(playing)
	changePlayerPosition(inputFile)
	// TODO make this work on Windows
	outputFile = sox.OpenWrite("default", inputFile.Signal(), nil, "alsa")
	if outputFile == nil {
		panic("Cannot open audio output devices")
	}
	//Now actually play
	chain = sox.CreateEffectsChain(inputFile.Encoding(), outputFile.Encoding())
	//make it output
	intermSignal := inputFile.Signal().Copy()
	//set input
	e := sox.CreateEffect(sox.FindEffect("input"))
	e.Options(inputFile)
	chain.Add(e, intermSignal, inputFile.Signal())
	e.Release()
	//set output
	e = sox.CreateEffect(sox.FindEffect("output"))
	e.Options(outputFile)
	chain.Add(e, intermSignal, inputFile.Signal())
	e.Release()
	//start the timer which keeps track of position in the file
	startTime = time.Now()
	lengthOfFile = getLengthOfFile(playing) //set length
	//process which also plays
	go chain.Flow()
	return chain, inputFile, outputFile
}

func fastForwardPlayer() {
	//save time and file
	if playerPosition == -1 {
		fmt.Println("Have to select a file to play to resume playback")
	} else {
		// TODO fix forward skip length
		// playerPosition += int(time.Since(startTime).Seconds()) + config.forwardSkipLength
		if playerPosition > int(lengthOfFile) {
			playerPosition = int(lengthOfFile) - 1
		}
		playerState = PlayerPlay
	}
}

func rewindPlayer() {
	//save time and file
	if playerPosition == -1 {
		fmt.Println("Have to select a file to play to resume playback")
	} else {
		// TODO fix rewind skip length
		// playerPosition += int(time.Since(startTime).Seconds()) - config.backwardSkipLength
		if playerPosition < 0 {
			playerPosition = 0
		}
		playerState = PlayerPlay
	}
}

func cleanupSoxData(chain **sox.EffectsChain, inputFile **sox.Format) {
	if inputFile != nil && *inputFile != nil {
		(*inputFile).Release()
		*inputFile = nil
	}
	if chain != nil && *chain != nil {
		(*chain).Release()
		*chain = nil
	}
}

func getLengthOfFile(fileName string) uint64 {
	inputFile := sox.OpenRead(fileName)
	seek := uint64(float64(inputFile.Signal().Length())/float64(inputFile.Signal().Channels())/float64(inputFile.Signal().Rate()) - 0.5)
	seek += seek % uint64(inputFile.Signal().Channels())
	return seek
}

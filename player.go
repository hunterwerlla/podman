package main

import (
	"fmt"
	"github.com/krig/go-sox"
	"time"
)

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

var (
	playerPosition int = -1
	startTime      time.Time
	lengthOfFile   uint64
	playerControl  chan PlayerState
	fileChannel    chan string
	exitChannel    chan bool
	playing        string
	playerState    PlayerState
)

// StartPlayer starts the global player. The player is global since there is only one of them
// I'm not really a fan of moving it into an object as it should not be reused
func StartPlayer() {
	playerState := make(chan PlayerState)
	fileChannel := make(chan string)
	exit := make(chan bool)
	go startPlayer(playerState, fileChannel, exit)
}

// DisposePlayer sends a signal to the player to destroy itself, and then waits for the player to exit
// This function will deadlock if called twice
func DisposePlayer() {
	playerControl <- ExitPlayer
	<-exitChannel
}

func PausePlayer() {
	playerPosition += int(time.Since(startTime).Seconds())
	playerState = Pause
}

func StopPlayer() {
	//reset position
	playerPosition = -1
	playerState = NothingPlaying
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
func startPlayer(playerState chan PlayerState, fileSelectionChannel chan string, exit chan bool) {
	//get rid of all stderr and stdout data
	//due to SOX outputting error messages
	// _, unused, _ := os.Pipe()
	// os.Stderr = unused
	// os.Stdout = unused
	playerControl = playerState
	fileChannel = fileSelectionChannel
	exitChannel = exit
	var (
		chain      *sox.EffectsChain
		inputFile  *sox.Format
		outputFile *sox.Format
		status     = NothingPlaying
		stopToExit = false
	)
	if !sox.Init() {
		panic("Unable to start the player")
	}
	defer sox.Quit()
	for stopToExit != true {
		status = NothingPlaying //reset status

		select {
		case status = <-playerControl:
		case playing = <-fileChannel:
			continue
		}

		switch status {
		case NothingPlaying:
			panic("invalid state when switching status, this should never happen")
		case Play:
			playFile(chain, inputFile, outputFile)
		case Pause: //case 1 pause
			//save time and file
			PausePlayer()
			cleanupSoxData(chain, inputFile, outputFile)
		case Stop:
			StopPlayer()
			cleanupSoxData(chain, inputFile, outputFile)
		case FastForward:
			fastForwardPlayer()
		case Rewind:
			rewindPlayer()
		case ExitPlayer:
			cleanupSoxData(chain, inputFile, outputFile)
			stopToExit = true
		}
	}
	exit <- true
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

func playFile(chain *sox.EffectsChain, inputFile *sox.Format, outputFile *sox.Format) {
	playerState = Play
	inputFile = sox.OpenRead(playing)
	changePlayerPosition(inputFile)
	//try two audio output methods
	outputFile = sox.OpenWrite("default", inputFile.Signal(), nil, "alsa")
	if outputFile == nil {
		outputFile = sox.OpenWrite("default", inputFile.Signal(), nil, "pulseaudio")
		if outputFile == nil {
			panic("Cannot open audio output devices")
		}
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
		playerState = Play
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
		playerState = Play
	}
}

func cleanupSoxData(chain *sox.EffectsChain, inputFile *sox.Format, outputFile *sox.Format) {
	if inputFile != nil {
		inputFile.Release()
	}
	if outputFile != nil {
		outputFile.Release()
	}
	if chain != nil {
		chain.Release()
	}
}

func getLengthOfFile(fileName string) uint64 {
	inputFile := sox.OpenRead(fileName)
	seek := uint64(float64(inputFile.Signal().Length())/float64(inputFile.Signal().Channels())/float64(inputFile.Signal().Rate()) - 0.5)
	seek += seek % uint64(inputFile.Signal().Channels())
	return seek
}

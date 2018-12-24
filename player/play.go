package main

import (
	"fmt"
	"github.com/krig/go-sox"
	"os"
	"time"
)

var (
	playerPosition int = -1
	startTime      time.Time
	lengthOfFile   uint64
)

//this runs on its own thread to start/stop and select the media that is playing
func play(exit chan bool) {
	//get rid of all stderr and stdout data
	//due to SOX outputting error messages
	_, unused, _ := os.Pipe()
	os.Stderr = unused
	os.Stdout = unused
	var (
		chain      *sox.EffectsChain = nil
		inputFile  *sox.Format       = nil
		outputFile *sox.Format       = nil
		status     PlayerState       = NothingPlaying
		stopToExit bool              = false
	)
	const (
		CLEAR = ""
		EMPTY = ""
	)
	if !sox.Init() {
		panic("Unable to start the player")
	}
	defer sox.Quit()
	for stopToExit != true {
		status = NothingPlaying //reset status
		status = <-globals.playerControl
		switch status {
		case NothingPlaying:
			panic("invalid state when switching status, this should never happen")
		case Play:
			startPlayer(chain, inputFile, outputFile)
		case Pause: //case 1 pause
			//save time and file
			pausePlayer()
			cleanupSoxData(chain, inputFile, outputFile)
		case Stop:
			stopPlayer()
			cleanupSoxData(chain, inputFile, outputFile)
		case FastForward:
			fastforward()
		case Rewind:
			rewind()
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

func startPlayer(chain *sox.EffectsChain, inputFile *sox.Format, outputFile *sox.Format) {
	inputFile = sox.OpenRead(globals.Playing)
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
	lengthOfFile = getLengthOfFile(globals.Playing) //set length
	//process which also plays
	go chain.Flow()
}
func pausePlayer() {
	playerPosition += int(time.Since(startTime).Seconds())
	globals.playerState = Pause
}

func stopPlayer() {
	//reset position
	playerPosition = -1
	globals.playerState = NothingPlaying
	globals.Playing = ""
	lengthOfFile = 0 //set length
}

func fastforward() {
	//save time and file
	if playerPosition == -1 {
		fmt.Println("Have to select a file to play to resume playback")
	} else {
		playerPosition += int(time.Since(startTime).Seconds()) + globals.Config.forwardSkipLength
		if playerPosition > int(lengthOfFile) {
			playerPosition = int(lengthOfFile) - 1
		}
		globals.playerState = Play
	}
}

func rewind() {
	//save time and file
	if playerPosition == -1 {
		fmt.Println("Have to select a file to play to resume playback")
	} else {
		playerPosition += int(time.Since(startTime).Seconds()) - globals.Config.backwardSkipLength
		if playerPosition < 0 {
			playerPosition = 0
		}
		globals.playerState = Play
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

func getLengthOfPlayingFile() uint64 {
	return lengthOfFile
}

func getPlayerState() int {
	return 0
}

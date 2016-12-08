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
		status     int               = _nothing
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
		status = _nothing //reset status
		status = <-globals.playerControl
		switch status {
		case _nothing:
			panic("invalid state when switching status, this should never happen")
		case _play:
			startPlayer(chain, inputFile, outputFile)
		case _pause: //case 1 pause
			//save time and file
			pausePlayer()
			cleanupSoxData(chain, inputFile, outputFile)
		case _stop:
			stopPlayer()
			cleanupSoxData(chain, inputFile, outputFile)
		case _ff:
			fastforward()
		case _rw:
			rewind()
		case _exit:
			cleanupSoxData(chain, inputFile, outputFile)
			stopToExit = true
		}
	}
	exit <- true
}

func changePlayerPosition(inputFile *sox.Format) {
	if globals.playerPosition < 0 {
		globals.playerPosition = 0
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
	interm_signal := inputFile.Signal().Copy()
	//set input
	e := sox.CreateEffect(sox.FindEffect("input"))
	e.Options(inputFile)
	chain.Add(e, interm_signal, inputFile.Signal())
	e.Release()
	//set output
	e = sox.CreateEffect(sox.FindEffect("output"))
	e.Options(outputFile)
	chain.Add(e, interm_signal, inputFile.Signal())
	e.Release()
	//start the timer which keeps track of position in the file
	startTime = time.Now()
	globals.LengthOfFile = getLengthOfFile(globals.Playing) //set length
	//process which also plays
	go chain.Flow()
}
func pausePlayer() {
	globals.playerPosition += int(time.Since(startTime).Seconds())
	globals.playerState = _pause
}
func stopPlayer() {
	//reset position
	globals.playerPosition = -1
	globals.playerState = _nothing
	globals.Playing = ""
	globals.LengthOfFile = 0 //set length
}

func fastforward() {
	//save time and file
	if playerPosition == -1 {
		fmt.Println("Have to select a file to play to resume playback")
	} else {
		playerPosition += int(time.Since(startTime).Seconds()) + globals.Config.forwardSkipLength
		if playerPosition > int(globals.LengthOfFile) {
			playerPosition = int(globals.LengthOfFile) - 1
		}
		globals.playerState = _play
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
		globals.playerState = _play
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

package main

import (
	"fmt"
	"github.com/krig/go-sox"
)

//this runs on its own thread to start/stop and select the media that is playing, it will also skip ahead in the future
//TODO make it skip ahead
func play(fileNameChannel chan string, control chan int) {
	var (
		chain   *sox.EffectsChain = nil
		inFile  *sox.Format       = nil
		outFile *sox.Format       = nil
	)
	if !sox.Init() {
		panic("Unable to start the player")
	}
	defer sox.Quit()
	for {
		//wait for a signal
		status := -1
		fileName := ""
		select {
		case fileName = <-fileNameChannel:
		case status = <-control:
		}
		//if filname is not empty, then new filename recieved
		if fileName != "" {
			//when switching stop before playing new file
			if chain != nil {
				chain.DeleteAll()
				chain.Release()
			}
			if inFile != nil {
				inFile.Release()
			}
			if outFile != nil {
				outFile.Release()
			}
			//stop sox if it's running
			//now start playing
			fmt.Println("playing")
			//first time get the fileName
			inFile = sox.OpenRead(fileName)
			//TODO make it work on windows too
			outFile = sox.OpenWrite("default", inFile.Signal(), nil, "alsa")
			if outFile == nil {
				outFile = sox.OpenWrite("default", inFile.Signal(), nil, "pulseaudio")
				if outFile == nil {
					panic("Cannot open audio output devices")
				}
			}
			//Now actually play
			chain = sox.CreateEffectsChain(inFile.Encoding(), outFile.Encoding())
			//make it output
			interm_signal := inFile.Signal().Copy()
			//set input
			e := sox.CreateEffect(sox.FindEffect("input"))
			e.Options(inFile)
			chain.Add(e, interm_signal, inFile.Signal())
			e.Release()
			//set output
			e = sox.CreateEffect(sox.FindEffect("output"))
			e.Options(outFile)
			chain.Add(e, interm_signal, inFile.Signal())
			e.Release()
			//process which also plays
			go chain.Flow()
		} else {
			switch status {
			case 1:
				return
			}
		}
	}
	if inFile != nil {
		inFile.Release()
	}
	if outFile != nil {
		outFile.Release()
	}
	if chain != nil {
		chain.Release()
	}
}

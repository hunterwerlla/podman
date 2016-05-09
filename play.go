package main

import (
	"fmt"
	"github.com/krig/go-sox"
	"time"
)

//this runs on its own thread to start/stop and select the media that is playing, it will also skip ahead in the future
//TODO make it skip ahead
//Control reference: 0 is play, 1 is pause, 2 is stop, 3 is skip ahead, 4 is reverse
func play(fileNameChannel chan string, control chan int) {
	var (
		chain     *sox.EffectsChain = nil
		inFile    *sox.Format       = nil
		outFile   *sox.Format       = nil
		position  int               = 0
		startTime time.Time
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
				chain.Release()
				chain = nil
			}
			if inFile != nil {
				inFile.Release()
				inFile = nil
			}
			if outFile != nil {
				outFile.Release()
				outFile = nil
			}
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
			//start the timer which keeps track of position in the file
			startTime = time.Now()
			//reset position
			position = 0
			//process which also plays
			go chain.Flow()
		} else {
			switch status {
			case 0: //case 0 play, only works after pause
				fmt.Printf("will start playing at %d\n", position)
				return
			case 1: //case 1 pause
				//save time and stop
				position = int(time.Since(startTime).Seconds())
				if chain != nil {
					chain.Release()
					chain = nil
				}
				if inFile != nil {
					inFile.Release()
					inFile = nil
				}
				if outFile != nil {
					outFile.Release()
					outFile = nil
				}
			case 2: //case 2 stop
				if chain != nil {
					chain.Release()
					chain = nil
				}
				if inFile != nil {
					inFile.Release()
					inFile = nil
				}
				if outFile != nil {
					outFile.Release()
					outFile = nil
				}
			case 3:
				return
			case 4:
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

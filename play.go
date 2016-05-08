package main

import (
	"github.com/krig/go-sox"
)

//this runs on its own thread to start/stop and select the media that is playing, it will also skip ahead in the future
//TODO make it skip ahead
func play(fileNameChannel chan string, control chan int) {
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
			if !sox.Init() {
				panic("Unable to start the player")
			}
			defer sox.Quit()
			//first time get the fileName
			inFile := sox.OpenRead(fileName)
			defer inFile.Release()
			//TODO make it work on windows too
			out := sox.OpenWrite("default", inFile.Signal(), nil, "alsa")
			if out == nil {
				out = sox.OpenWrite("default", inFile.Signal(), nil, "pulseaudio")
				if out == nil {
					panic("Cannot open audio output devices")
				}
			}
			defer out.Release()
			//Now actually play
			chain := sox.CreateEffectsChain(inFile.Encoding(), out.Encoding())
			defer chain.Release()
			//make it output
			interm_signal := inFile.Signal().Copy()
			//set input
			e := sox.CreateEffect(sox.FindEffect("input"))
			e.Options(inFile)
			chain.Add(e, interm_signal, inFile.Signal())
			e.Release()
			//set output
			e = sox.CreateEffect(sox.FindEffect("output"))
			e.Options(out)
			chain.Add(e, interm_signal, inFile.Signal())
			e.Release()
			//process which also plays
			chain.Flow()
		} else {
			switch status {
			case 1:
				return
			}
		}
	}
}

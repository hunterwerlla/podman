package main

import (
	"fmt"
	ui "github.com/gizak/termui"
)

const (
	podmanHeader = "" +
		" _____         _                  \n" +
		"|  _  | ___  _| | _____  ___  ___ \n" +
		"|   __|| . || . ||     || .'||   |\n" +
		"|__|   |___||___||_|_|_||__,||_|_|\n"
)

func produceHeader(width int) *ui.Paragraph {
	headerWidget := ui.NewParagraph(podmanHeader)
	headerWidget.Height = 5
	headerWidget.Width = width
	headerWidget.TextFgColor = ui.ColorBlack
	headerWidget.BorderTop = false
	headerWidget.BorderLeft = false
	headerWidget.BorderRight = false
	headerWidget.BorderBottom = true
	headerWidget.WrapLength = 0
	return headerWidget
}

func handleKeyboard(configuration *Configuration, event ui.Event) {
	if event.ID == configuration.PlayKeybind {
		fmt.Println("play")
	}
}

func StartTui(configuration *Configuration) {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	width := ui.Body.Width
	headerWidget := produceHeader(width)
	ui.Render(headerWidget)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			if e.ID == "<C-c>" {
				break
			} else {
				handleKeyboard(configuration, e)
			}
		}
	}
}

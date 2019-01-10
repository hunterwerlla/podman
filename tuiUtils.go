package main

import (
	ui "github.com/gizak/termui"
)

func getForegroundColorForTheme(configuration *Configuration) ui.Attribute {
	if configuration.Theme == ThemeLight {
		return ui.ColorBlack
	} else if configuration.Theme == ThemeDark {
		return ui.ColorWhite
	}
	return ui.ColorBlack
}

func getForegroundColorForThemeString(configuration *Configuration) string {
	if configuration.Theme == ThemeLight {
		return "black"
	} else if configuration.Theme == ThemeDark {
		return "white"
	}
	return "black"
}

func getBackgroundColorForTheme(configuration *Configuration) ui.Attribute {
	if configuration.Theme == ThemeLight {
		return ui.ColorWhite
	} else if configuration.Theme == ThemeDark {
		return ui.ColorBlack
	}
	return ui.ColorWhite
}

func getBackgroundColorForThemeString(configuration *Configuration) string {
	if configuration.Theme == ThemeLight {
		return "white"
	} else if configuration.Theme == ThemeDark {
		return "black"
	}
	return "white"
}

func termuiStyleText(text string, fgcolor string, bgcolor string) string {
	text = "[" + text + "](fg-" + fgcolor + ",bg-" + string(bgcolor) + ")"
	return text
}

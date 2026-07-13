package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const spicetifyColorTemplate = `[pywal]
text = c0c0c0
subtext = a0a0a0
nav-active-text = 101010
main = 101010
sidebar = 101010
player = 101010
card = 101010
shadow = 101010
main-secondary = 101010
button = 6f8faf
button-secondary = c0c0c0
button-active = 6f8faf
button-disabled = 353535
nav-active = 8aa0c8
play-button = 6f8faf
tab-active = 101010
notification = c0c0c0
notification-error = 353535
playback-bar = 6f8faf
misc = c0c0c0
`

func UpdateSpicetify(pywalColors map[string]string, homeDir string) error {
	themeDir := filepath.Join(homeDir, ".config", "spicetify", "Themes", "Sleek")
	pathToTheme := filepath.Join(themeDir, "color.ini")
	configBytes, err := os.ReadFile(pathToTheme)
	if err != nil {
		configBytes, err = os.ReadFile(filepath.Join(themeDir, "color.ini.template"))
		if err != nil {
			configBytes = []byte(spicetifyColorTemplate)
		}
	}

	configAsArray := strings.Split(string(configBytes), "\n")

	for i, line := range configAsArray {
		if !strings.Contains(line, "=") {
			continue
		}

		splitLine := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(splitLine[0])

		var newColorKey string

		switch key {
		case "main", "sidebar", "nav-active-text", "main-secondary", "player", "card", "shadow", "tab-active":
			newColorKey = "background"

		case "subtext":
			newColorKey = "foreground"

		case "playback-bar", "play-button", "button", "button-active":
			newColorKey = "color4"

		case "text", "button-secondary", "notification", "misc":
			newColorKey = "color6"

		case "nav-active":
			newColorKey = "color5"
		default:
			continue
		}

		if newColor, ok := pywalColors[newColorKey]; ok {
			cleanColor := strings.Trim(newColor, "'# ")
			newLine := fmt.Sprintf("%s = %s", key, cleanColor)
			configAsArray[i] = newLine
		}
	}

	output := strings.Join(configAsArray, "\n")
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(pathToTheme, []byte(output), 0644)
}

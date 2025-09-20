package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func UpdateSpicetify(pywalColors map[string]string, homeDir string) error {
	pathToTheme := homeDir + "/.config/spicetify/Themes/Sleek/color.ini"
	file, err := os.Open(pathToTheme)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)

	var configAsArray []string
	for scanner.Scan() {
		configAsArray = append(configAsArray, scanner.Text())
	}

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

		case "text", "button-secondary", "notfication", "misc":

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
	err = os.WriteFile(pathToTheme, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

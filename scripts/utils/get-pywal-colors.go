package utils

import (
	"bufio"
	"os"
	"strings"
)

func GetPywalColors() (map[string]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	pywalColors := make(map[string]string)

	pywalColorsPath := homeDir + "/.cache/wal/colors.sh"
	pywalFile, err := os.Open(pywalColorsPath)
	if err != nil {
		return nil, err
	}
	defer pywalFile.Close()

	scanner := bufio.NewScanner(pywalFile)

	for scanner.Scan() {
		split := strings.Split(scanner.Text(), "=")
		if len(split) > 1 {
			pywalColors[split[0]] = split[1]
		}
	}
	return pywalColors, nil
}

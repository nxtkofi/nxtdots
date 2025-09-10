package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func returnOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	homeDir, err := os.UserHomeDir()
	returnOnErr(err)
	cavaConfigPath := homeDir + "/.config/cava/config"
	pywalColors, error := os.Open(homeDir + "/.cache/wal/colors")
	returnOnErr(error)
	scanner := bufio.NewScanner(pywalColors)

	var allPywalColors []string

	for scanner.Scan() {
		allPywalColors = append(allPywalColors, scanner.Text())
	}

	var newColors []string
	if len(allPywalColors) >= 3 {
		newColors = allPywalColors[0:4]
	}

	fmt.Println(newColors)
	cavaConfig, error := os.Open(cavaConfigPath)
	returnOnErr(error)

	scanner = bufio.NewScanner(cavaConfig)
	var cavaConfigAsArray []string
	for scanner.Scan() {
		cavaConfigAsArray = append(cavaConfigAsArray, scanner.Text())
	}

	for i, line := range cavaConfigAsArray {
		if strings.HasPrefix(line, "gradient_color_1") {
			cavaConfigAsArray[i] = "gradient_color_1 = '" + newColors[1] + "'"
		}

		if strings.HasPrefix(line, "gradient_color_2") {
			cavaConfigAsArray[i] = "gradient_color_2 = '" + newColors[2] + "'"
		}

		if strings.HasPrefix(line, "gradient_color_3") {
			cavaConfigAsArray[i] = "gradient_color_3 = '" + newColors[3] + "'"
			break // Found last gradient, no need to read the rest of the file :)
		}
	}

	output := strings.Join(cavaConfigAsArray, "\n")
	err = os.WriteFile(cavaConfigPath, []byte(output), 0644)
	returnOnErr(err)
}

package utils

import (
	"bufio"
	"os"
	"strings"
)

func UpdateCavaGradient(cavaConfigPath string, allPywalColors []string) error {
	var newColors []string
	if len(allPywalColors) >= 4 {
		newColors = allPywalColors[0:4]
	} else {
		return nil
	}

	cavaConfig, err := os.Open(cavaConfigPath)
	if err != nil {
		return err
	}
	defer cavaConfig.Close()

	scanner := bufio.NewScanner(cavaConfig)
	var cavaConfigAsArray []string
	for scanner.Scan() {
		cavaConfigAsArray = append(cavaConfigAsArray, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
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
		}
	}

	output := strings.Join(cavaConfigAsArray, "\n")
	err = os.WriteFile(cavaConfigPath, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

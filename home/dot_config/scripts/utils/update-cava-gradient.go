package utils

import (
	"bufio"
	"os"
	"strings"
)

func UpdateCavaGradient(cavaConfigPath string, allPywalColors map[string]string) error {
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
			cavaConfigAsArray[i] = "gradient_color_1 = " + allPywalColors["color1"] + ""
		}

		if strings.HasPrefix(line, "gradient_color_2") {
			cavaConfigAsArray[i] = "gradient_color_2 = " + allPywalColors["color2"] + ""
		}

		if strings.HasPrefix(line, "gradient_color_3") {
			cavaConfigAsArray[i] = "gradient_color_3 = " + allPywalColors["color3"] + ""
		}
	}

	output := strings.Join(cavaConfigAsArray, "\n")
	err = os.WriteFile(cavaConfigPath, []byte(output), 0644)
	if err != nil {
		return err
	}

	return nil
}

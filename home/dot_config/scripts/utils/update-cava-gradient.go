package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func UpdateCavaGradient(homeDir string, allPywalColors map[string]string) error {
	themeDir := filepath.Join(homeDir, ".config", "cava", "themes")
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		return err
	}

	theme := fmt.Sprintf(`gradient = 1
gradient_count = 3
gradient_color_1 = %s
gradient_color_2 = %s
gradient_color_3 = %s
`, allPywalColors["color1"], allPywalColors["color2"], allPywalColors["color3"])

	return os.WriteFile(filepath.Join(themeDir, "pywal.generated.local"), []byte(theme), 0o644)
}

package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func UpdateCavaGradient(homeDir string, allPywalColors map[string]string) error {
	themeDir := filepath.Join(homeDir, ".config", "cava", "themes")
	if err := os.MkdirAll(themeDir, 0o755); err != nil {
		return err
	}

	theme := fmt.Sprintf(`[color]
gradient = 1
gradient_count = 3
gradient_color_1 = '%s'
gradient_color_2 = '%s'
gradient_color_3 = '%s'
`, strings.Trim(allPywalColors["color1"], "'"), strings.Trim(allPywalColors["color2"], "'"), strings.Trim(allPywalColors["color3"], "'"))

	if err := os.WriteFile(filepath.Join(themeDir, "pywal.generated.local"), []byte(theme), 0o644); err != nil {
		return err
	}

	if err := exec.Command("pkill", "-USR2", "-x", "cava").Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
			return nil
		}
		return fmt.Errorf("reload cava colors: %w", err)
	}

	return nil
}

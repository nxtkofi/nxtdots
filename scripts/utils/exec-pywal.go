package utils

import (
	"os/exec"
)

func ExecPywal(wallpaperFilePath string) error {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	if string(output) == "'prefer-dark'" {
		cmd := exec.Command("wal", "-q", "-i", wallpaperFilePath)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	} else if string(output) == "'prefer-light'" {
		cmd := exec.Command("wal", "-q", "-l", "-i", wallpaperFilePath)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}
	return nil
}

package utils

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExecPywal(wallpaperFilePath string) error {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	output, err := cmd.Output()
	if err != nil {
		return err
	}
	colorscheme := strings.TrimSpace(string(output))

	if colorscheme == "'prefer-dark'" {
		fmt.Printf("prefer-darky")
		cmd := exec.Command("wal", "-q", "-i", wallpaperFilePath)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	} else if colorscheme == "'prefer-light'" {
		fmt.Printf("prefer-ligthy")
		cmd := exec.Command("wal", "-q", "-l", "-i", wallpaperFilePath)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}
	return resetKittyIfItsRunning()
}

func resetKittyIfItsRunning() error {
	cmd := exec.Command("pgrep", "-x", "kitty")
	err := cmd.Run()
	if err != nil {
		cmd := exec.Command("pkill", "-USR1", "-x", "kitty")
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

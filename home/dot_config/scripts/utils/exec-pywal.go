package utils

import (
	"os/exec"
)

func ExecPywal(colorMode XDGColorScheme, wallpaperFilePath string) error {

	switch colorMode {
	case "'prefer-dark'":
		cmd := exec.Command("wal", "-q", "-i", wallpaperFilePath)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	case "'prefer-light'":
		cmd := exec.Command("wal", "-q", "-l", "-i", wallpaperFilePath)
		_, err := cmd.Output()
		if err != nil {
			return err
		}
	}
	resetKittyIfItsRunning()
	return nil
}

func resetKittyIfItsRunning() {
	cmd := exec.Command("pkill", "-SIGUSR1", "kitty")
	// Ignore error - pkill returns exit status 1 if no process is found, which is expected
	cmd.Run()
}

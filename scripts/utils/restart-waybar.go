package utils

import (
	"os/exec"
)

func RestartWaybar() error {
	killCmd := exec.Command("killall", "waybar")
	_ = killCmd.Run() 

	startCmd := exec.Command("waybar")
	err := startCmd.Start()
	if err != nil {
		return err
	}

	return nil
}

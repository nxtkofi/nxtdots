package utils

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

func Install() {
	homedir, err := os.UserHomeDir()
	ReturnOnErr(err)
	yayCmd := exec.Command("yay", "-S", "fzf", "waybar", "downgrade", "vesktop", "walcord", "spicetify-cli", "python-pywal16", "magick", "nmtui", "bluetuith", "power-profiles-daemon", "zen-browser-bin", "hyprland", "--no-confirm")
	copyFromUsrShareToLocalAndPerformOverwrite("org.moson.pacseek.desktop", "Exec", "Exec=kitty --class Pacseek pacseek")
	copyFromUsrShareToLocalAndPerformOverwrite("spotify.desktop", "Exec", "Exec="+homedir+"/.config/settings/launch-spotify.sh")
	err = yayCmd.Run()
	ReturnOnErr(err)

	chmodRecursiveCmd := exec.Command("sudo", "chmod", "a+wr", "/opt/spotify/Apps -R")
	err = chmodRecursiveCmd.Run()
	ReturnOnErr(err)

	chmodCmd := exec.Command("sudo", "chmod", "a+wr", "/opt/spotify")
	err = chmodCmd.Run()
	ReturnOnErr(err)

	powerProfileDaemonEnable := exec.Command("sudo", "systemctl", "enable", "--now", "power-profiles-daemon")
	err = powerProfileDaemonEnable.Run()
	ReturnOnErr(err)

	spicetifyConfig := exec.Command("spicetify", "config", "current_theme", "Sleek")
	err = spicetifyConfig.Run()
	ReturnOnErr(err)

	spicetifyApply := exec.Command("spicetify", "apply")
	err = spicetifyApply.Run()
	ReturnOnErr(err)

	runWaybar := exec.Command("waybar")
	err = runWaybar.Run()
	ReturnOnErr(err)
}

func copyFromUsrShareToLocalAndPerformOverwrite(ogFileName, keyString, replaceValue string) {
	homedir, err := os.UserHomeDir()
	ReturnOnErr(err)
	copyFile, err := os.Open("/usr/share/applications/" + ogFileName)
	ReturnOnErr(err)
	newFile, err := os.Open(homedir + "/.local/share/applications/" + ogFileName)
	ReturnOnErr(err)
	defer copyFile.Close()
	defer newFile.Close()
	scanner := bufio.NewScanner(newFile)
	var fileContent []string

	for scanner.Scan() {
		fileContent = append(fileContent, scanner.Text())
	}

	for i, line := range fileContent {
		var split []string
		split = strings.Split(line, "=")
		if strings.TrimSpace(split[0]) == keyString {
			fileContent[i] = replaceValue
		}
	}

	err = os.WriteFile(newFile.Name(), []byte(strings.Join(fileContent, "\n")), 0644)
	ReturnOnErr(err)
}

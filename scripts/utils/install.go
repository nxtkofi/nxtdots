package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func Install() {
	homedir, err := os.UserHomeDir()
	ReturnOnErr(err)

	packages := []string{"fzf", "waybar", "downgrade", "vesktop", "walcord", "spicetify-cli", "python-pywal16", "magick", "nmtui", "bluetuith", "power-profiles-daemon", "zen-browser-bin", "hyprland", "spotify"}

	if _, err := exec.LookPath("yay"); err == nil {
		fmt.Println("Using yay for package installation...")
		cmd := exec.Command("yay", append([]string{"-S", "--no-confirm"}, packages...)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
	} else {
		fmt.Println("yay not found, using pacman for official packages...")
		officialPackages := []string{"fzf", "waybar", "downgrade", "magick", "nmtui", "bluetuith", "power-profiles-daemon"}
		args := append([]string{"pacman", "-S", "--noconfirm"}, officialPackages...)
		cmd := exec.Command("sudo", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
		if err == nil {
			fmt.Println("AUR packages (vesktop, walcord, spicetify-cli, python-pywal16, zen-browser-bin, hyprland) need to be installed manually or with an AUR helper like yay")
		}
	}
	ReturnOnErr(err)

	copyFromUsrShareToLocalAndPerformOverwrite("org.moson.pacseek.desktop", "Exec", "Exec=kitty --class Pacseek pacseek")
	copyFromUsrShareToLocalAndPerformOverwrite("spotify.desktop", "Exec", "Exec="+homedir+"/.config/settings/launch-spotify.sh")

	chmodRecursiveCmd := exec.Command("sudo", "chmod", "-R", "a+wr", "/opt/spotify/Apps")
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

	fmt.Println("Installation complete. You can now run 'waybar' to start waybar.")
}

func copyFromUsrShareToLocalAndPerformOverwrite(ogFileName, keyString, replaceValue string) {
	homedir, err := os.UserHomeDir()
	ReturnOnErr(err)

	// Ensure local applications directory exists
	localAppsDir := homedir + "/.local/share/applications/"
	err = os.MkdirAll(localAppsDir, 0755)
	ReturnOnErr(err)

	// Copy file from /usr/share/applications to local
	srcFile, err := os.Open("/usr/share/applications/" + ogFileName)
	ReturnOnErr(err)
	defer srcFile.Close()

	dstFile, err := os.Create(localAppsDir + ogFileName)
	ReturnOnErr(err)
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	ReturnOnErr(err)
	srcFile.Close()
	dstFile.Close()

	newFile, err := os.Open(localAppsDir + ogFileName)
	ReturnOnErr(err)
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

	err = os.WriteFile(localAppsDir+ogFileName, []byte(strings.Join(fileContent, "\n")), 0644)
	ReturnOnErr(err)
}

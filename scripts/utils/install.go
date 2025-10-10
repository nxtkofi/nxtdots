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

	// Setup wallpapers and bashrc before package installation
	setupWallpapers(homedir)
	setupBashrc(homedir)

	packages := []string{"kitty", "fzf", "waybar", "downgrade", "vesktop", "walcord", "spicetify-cli", "python-pywal16", "imagemagick", "bluetui", "power-profiles-daemon", "zen-browser-bin", "hyprland", "spotify", "pacseek", "waypaper", "rofi", "hyprlock", "hyprpaper", "nautilus", "fastfetch", "starship", "zoxide", "noto-fonts-emoji", "ttf-jetbrains-mono-nerd", "ttf-firacode-nerd", "nerd-fonts-fira-code", "swaync", "xdg-desktop-portal", "xdg-desktop-portal-gtk", "xdg-desktop-portal-hyprland", "sddm", "qt6-svg", "qt6-virtualkeyboard", "qt6-multimedia-ffmpeg", "nvm"}

	if _, err := exec.LookPath("yay"); err == nil {
		fmt.Println("Using yay for package installation...")
		cmd := exec.Command("yay", append([]string{"-S", "--noconfirm"}, packages...)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		err = cmd.Run()
	} else {
		fmt.Println("yay not found, attempting to install yay...")
		if installYay() {
			fmt.Println("yay installed successfully, installing all packages...")
			cmd := exec.Command("yay", append([]string{"-S", "--noconfirm"}, packages...)...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			err = cmd.Run()
		} else {
			fmt.Println("Failed to install yay. Cannot proceed with package installation.")
			fmt.Println("Please install yay manually and run this script again.")
			return
		}
	}
	ReturnOnErr(err)

	copyFromUsrShareToLocalAndPerformOverwrite("org.moson.pacseek.desktop", "Exec", "Exec=kitty --class Pacseek pacseek")

	// Check if launch-spotify.sh exists before setting it
	launchSpotifyPath := homedir + "/.config/settings/launch-spotify.sh"
	if _, err := os.Stat(launchSpotifyPath); os.IsNotExist(err) {
		fmt.Printf("Warning: %s does not exist, using default spotify command\n", launchSpotifyPath)
		copyFromUsrShareToLocalAndPerformOverwrite("spotify.desktop", "Exec", "Exec=spotify")
	} else {
		copyFromUsrShareToLocalAndPerformOverwrite("spotify.desktop", "Exec", "Exec="+launchSpotifyPath)
	}

	chmodRecursiveCmd := exec.Command("sudo", "chmod", "-R", "a+wr", "/opt/spotify/Apps")
	err = chmodRecursiveCmd.Run()
	ReturnOnErr(err)

	chmodCmd := exec.Command("sudo", "chmod", "a+wr", "/opt/spotify")
	err = chmodCmd.Run()
	ReturnOnErr(err)

	powerProfileDaemonEnable := exec.Command("sudo", "systemctl", "enable", "--now", "power-profiles-daemon")
	err = powerProfileDaemonEnable.Run()
	ReturnOnErr(err)

	spicetifyBackupCreate := exec.Command("spicetify", "backup", "apply")
	err = spicetifyBackupCreate.Run()
	ReturnOnErr(err)

	spicetifyConfig := exec.Command("spicetify", "config", "current_theme", "Sleek")
	err = spicetifyConfig.Run()
	ReturnOnErr(err)

	spicetifyApply := exec.Command("spicetify", "apply")
	err = spicetifyApply.Run()
	ReturnOnErr(err)

	installSddmTheme()
	setupOneTimeCommands()
	updateWaybarThemePaths(homedir)

	fmt.Println("Installation complete. You can now run 'waybar' to start waybar.")
}

func copyFromUsrShareToLocalAndPerformOverwrite(ogFileName, keyString, replaceValue string) {
	homedir, err := os.UserHomeDir()
	ReturnOnErr(err)

	// Ensure local applications directory exists
	localAppsDir := homedir + "/.local/share/applications/"
	err = os.MkdirAll(localAppsDir, 0755)
	ReturnOnErr(err)

	// Check if source file exists
	srcPath := "/usr/share/applications/" + ogFileName
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		fmt.Printf("Warning: Source file %s does not exist, skipping...\n", srcPath)
		return
	}

	// Copy file from /usr/share/applications to local
	srcFile, err := os.Open(srcPath)
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
	ReturnOnErr(scanner.Err())

	for i, line := range fileContent {
		var split []string
		split = strings.Split(line, "=")
		if strings.TrimSpace(split[0]) == keyString {
			fileContent[i] = replaceValue
		}
	}

	err = os.WriteFile(localAppsDir+ogFileName, []byte(strings.Join(fileContent, "\n")+"\n"), 0644)
	ReturnOnErr(err)
}

func installYay() bool {
	if os.Geteuid() == 0 {
		fmt.Println("Cannot install yay as root. Need non-root user.")
		return false
	}

	fmt.Println("Installing git and base-devel...")
	cmd := exec.Command("sudo", "pacman", "-S", "--noconfirm", "git", "base-devel")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to install prerequisites: %v\n", err)
		return false
	}

	fmt.Println("Cloning yay from AUR...")
	cmd = exec.Command("git", "clone", "https://aur.archlinux.org/yay.git", "/tmp/yay-install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to clone yay: %v\n", err)
		exec.Command("rm", "-rf", "/tmp/yay-install").Run()
		return false
	}

	fmt.Println("Building and installing yay...")
	cmd = exec.Command("makepkg", "-si", "--noconfirm")
	cmd.Dir = "/tmp/yay-install"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Failed to build yay: %v\n", err)
		exec.Command("rm", "-rf", "/tmp/yay-install").Run()
		return false
	}

	// Clean up
	exec.Command("rm", "-rf", "/tmp/yay-install").Run()

	return true
}

func installSddmTheme() {
	fmt.Println("Installing SDDM Astronaut Theme...")

	// Clone the theme repository
	cloneTheme := exec.Command("sudo", "git", "clone", "-b", "master", "--depth", "1",
		"https://github.com/keyitdev/sddm-astronaut-theme.git",
		"/usr/share/sddm/themes/sddm-astronaut-theme")
	err := cloneTheme.Run()
	ReturnOnErr(err)

	// Copy fonts
	copyFonts := exec.Command("sudo", "cp", "-r",
		"/usr/share/sddm/themes/sddm-astronaut-theme/Fonts/",
		"/usr/share/fonts/")
	err = copyFonts.Run()
	ReturnOnErr(err)

	// Configure SDDM theme
	sddmConfig := `[Theme]
Current=sddm-astronaut-theme`
	err = os.WriteFile("/tmp/sddm.conf", []byte(sddmConfig), 0644)
	ReturnOnErr(err)

	copySddmConfig := exec.Command("sudo", "cp", "/tmp/sddm.conf", "/etc/sddm.conf")
	err = copySddmConfig.Run()
	ReturnOnErr(err)

	// Configure virtual keyboard
	virtKbdConfig := `[General]
InputMethod=qtvirtualkeyboard`
	err = os.WriteFile("/tmp/virtualkbd.conf", []byte(virtKbdConfig), 0644)
	ReturnOnErr(err)

	// Ensure sddm.conf.d directory exists
	mkdirSddm := exec.Command("sudo", "mkdir", "-p", "/etc/sddm.conf.d")
	err = mkdirSddm.Run()
	ReturnOnErr(err)

	copyVirtKbd := exec.Command("sudo", "cp", "/tmp/virtualkbd.conf", "/etc/sddm.conf.d/virtualkbd.conf")
	err = copyVirtKbd.Run()
	ReturnOnErr(err)

	// Enable SDDM service
	enableSddm := exec.Command("sudo", "systemctl", "enable", "sddm")
	err = enableSddm.Run()
	ReturnOnErr(err)

	// Clean up temp files
	os.Remove("/tmp/sddm.conf")
	os.Remove("/tmp/virtualkbd.conf")

	fmt.Println("SDDM Astronaut Theme has been installed and configured.")
}

func setupWallpapers(homedir string) {
	fmt.Println("Setting up wallpapers...")

	// Get current working directory (where the script is being run from)
	wd, err := os.Getwd()
	ReturnOnErr(err)

	assetsWallpaperPath := wd + "/assets/wallpaper"
	homeWallpaperPath := homedir + "/wallpaper"

	// Check if source assets/wallpaper directory exists
	if _, err := os.Stat(assetsWallpaperPath); os.IsNotExist(err) {
		fmt.Printf("Warning: Assets wallpaper directory %s does not exist, skipping wallpaper setup\n", assetsWallpaperPath)
		return
	}

	// Copy wallpaper directory from assets to home
	copyWallpapers := exec.Command("cp", "-r", assetsWallpaperPath, homeWallpaperPath)
	err = copyWallpapers.Run()
	ReturnOnErr(err)

	fmt.Printf("Wallpapers copied from %s to %s\n", assetsWallpaperPath, homeWallpaperPath)
}

func setupBashrc(homedir string) {
	fmt.Println("Setting up .bashrc...")

	bashrcContent := `for f in ~/.config/bashrc/*; do
if [ ! -d $f ] ;then
c=` + "`echo $f | sed -e \"s=.config/bashrc=.config/bashrc/custom=\"`" + `
[[ -f $c ]] && source $c || source $f
fi
done`

	bashrcPath := homedir + "/.bashrc"
	err := os.WriteFile(bashrcPath, []byte(bashrcContent), 0644)
	ReturnOnErr(err)

	fmt.Printf(".bashrc created at %s\n", bashrcPath)
}

func setupOneTimeCommands() {
	fmt.Println("Running one-time setup commands...")

	// Set initial XDG color scheme to prefers-dark
	setColorScheme := exec.Command("gsettings", "set", "org.gnome.desktop.interface", "color-scheme", "prefer-dark")
	err := setColorScheme.Run()
	ReturnOnErr(err)

	// Enable SDDM service
	enableSddm := exec.Command("sudo", "systemctl", "enable", "sddm")
	err = enableSddm.Run()
	ReturnOnErr(err)

	// Enable NetworkManager service
	enableNetworkManager := exec.Command("sudo", "systemctl", "enable", "NetworkManager")
	err = enableNetworkManager.Run()
	ReturnOnErr(err)

	fmt.Println("One-time setup commands completed")
}

func updateWaybarThemePaths(homedir string) {
	fmt.Println("Updating waybar theme paths...")

	// Update dark theme
	darkThemePath := homedir + "/.config/waybar/themes/nxtdots-pywal-dark.css"
	updateCSSFilePath(darkThemePath, homedir)

	// Update light theme
	lightThemePath := homedir + "/.config/waybar/themes/nxtdots-pywal-light.css"
	updateCSSFilePath(lightThemePath, homedir)

	fmt.Println("Waybar theme paths updated")
}

func updateCSSFilePath(filePath, homedir string) {
	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Warning: Could not read %s, skipping path update\n", filePath)
		return
	}

	// Replace hardcoded path with actual home directory
	oldPath := "file:///home/nxtdots/.cache/wal/colors-waybar.css"
	newPath := "file://" + homedir + "/.cache/wal/colors-waybar.css"
	updatedContent := strings.ReplaceAll(string(content), oldPath, newPath)

	// Write back to file
	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	ReturnOnErr(err)

	fmt.Printf("Updated path in %s\n", filePath)
}

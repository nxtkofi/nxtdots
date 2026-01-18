package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func Install() {
	err := InitLogger()
	if err != nil {
		fmt.Printf("Warning: Could not initialize logging: %v\n", err)
	}
	defer CloseLogger()

	homedir, err := os.UserHomeDir()
	ReturnOnErr(err)

	LogInfo("Starting installation process")
	LogInfo(fmt.Sprintf("Home directory: %s", homedir))

	setupWallpapers(homedir)
	setupBashrc(homedir)

	if _, err := exec.LookPath("yay"); err != nil {
		LogInfo("yay not found, attempting to install yay")
		if !installYay() {
			LogError("Failed to install yay. Cannot proceed with package installation", nil)
			LogInfo("Please install yay manually and run this script again")
			return
		}
		LogInfo("yay installed successfully")
	}

	corePackages := []string{
		"hyprland", "waybar", "kitty", "rofi", "wofi", "wofi-emoji",
		"swaync", "hyprlock", "hypridle", "hyprpaper", "hyprshot", "hyprsunset",
		"xdg-desktop-portal", "xdg-desktop-portal-gtk", "xdg-desktop-portal-hyprland",
		"sddm", "qt6-svg", "qt6-virtualkeyboard", "qt6-multimedia-ffmpeg",
	}

	aurPackages := []string{
		"vesktop", "walcord", "spicetify-cli", "zen-browser-bin",
		"spotify", "pacseek", "swww",
	}

	utilityPackages := []string{
		"fzf", "ripgrep", "nvim", "wl-clipboard", "cliphist", "jq",
		"brightnessctl", "imagemagick", "bluetui", "power-profiles-daemon",
		"downgrade", "python-pywal16", "missioncenter", "nautilus",
		"fastfetch", "starship", "zoxide", "bash-completion", "nvm",
		"eza", "mpv", "papirus-icon-theme", "wtype",
	}

	fontPackages := []string{
		"noto-fonts-emoji", "ttf-jetbrains-mono-nerd", "ttf-firacode-nerd",
		"nerd-fonts-fira-code", "ttf-0xproto-nerd",
	}

	LogInfo("Installing core Hyprland packages")
	installPackageGroup(corePackages)

	LogInfo("Installing AUR packages")
	installPackageGroup(aurPackages)

	LogInfo("Installing utility packages")
	installPackageGroup(utilityPackages)

	LogInfo("Installing fonts")
	installPackageGroup(fontPackages)

	LogInfo("Validating critical packages installation")
	criticalPackages := []string{"hyprland", "waybar", "kitty", "nvim"}
	for _, pkg := range criticalPackages {
		if !isPackageInstalled(pkg) {
			LogError(fmt.Sprintf("Critical package '%s' is not installed!", pkg), nil)
			fmt.Printf("Installation may have failed. Please check logs and install missing packages manually.\n")
		}
	}

	LogInfo("Updating Waybar theme paths")
	updateWaybarThemePaths(homedir)
	updateWaybarMainTheme(homedir)

	LogInfo("Enabling power-profiles-daemon")
	powerProfileDaemonEnable := exec.Command("sudo", "systemctl", "enable", "--now", "power-profiles-daemon")
	LogCommand("sudo systemctl enable --now power-profiles-daemon")
	err = powerProfileDaemonEnable.Run()
	ReturnOnErr(err)

	LogInfo("Installing SDDM theme")
	installSddmTheme()
	LogInfo("Running one-time setup commands")
	setupOneTimeCommands()

	LogInfo("Setting up desktop entries")
	setupDesktopEntries(homedir)

	LogInfo("Checking for NVIDIA GPU")
	setupNvidiaIfDetected(homedir)

	LogInfo("Running initial pywal setup")
	runInitialPywalSetup(homedir)

	LogInfo("Installation complete. You can now run 'waybar' to start waybar")
	fmt.Println("Installation complete. You can now run 'waybar' to start waybar.")
}

func RiceSpotify() {
	err := InitLogger()
	if err != nil {
		fmt.Printf("Warning: Could not initialize logging: %v\n", err)
	}
	defer CloseLogger()

	homedir, err := os.UserHomeDir()
	ReturnOnErr(err)

	LogInfo("Starting Spotify ricing process")
	LogInfo(fmt.Sprintf("Home directory: %s", homedir))

	launchSpotifyPath := homedir + "/.config/settings/launch-spotify.sh"
	if _, err := os.Stat(launchSpotifyPath); os.IsNotExist(err) {
		fmt.Printf("Warning: %s does not exist, using default spotify command\n", launchSpotifyPath)
		copyFromUsrShareToLocalAndPerformOverwrite("spotify.desktop", "Exec", "Exec=spotify")
	} else {
		copyFromUsrShareToLocalAndPerformOverwrite("spotify.desktop", "Exec", "Exec="+launchSpotifyPath)
	}

	LogInfo("Setting Spotify permissions")
	chmodRecursiveCmd := exec.Command("sudo", "chmod", "-R", "a+wr", "/opt/spotify/Apps")
	LogCommand("sudo chmod -R a+wr /opt/spotify/Apps")
	err = chmodRecursiveCmd.Run()
	ReturnOnErr(err)

	chmodCmd := exec.Command("sudo", "chmod", "a+wr", "/opt/spotify")
	LogCommand("sudo chmod a+wr /opt/spotify")
	err = chmodCmd.Run()
	ReturnOnErr(err)

	if _, err := os.Stat("/opt/spotify"); os.IsNotExist(err) {
		LogError("Spotify not found, skipping spicetify configuration", err)
		return
	}

	LogInfo("Configuring Spicetify")
	spicetifyInitCmd := exec.Command("spicetify")
	LogCommand("spicetify")
	spicetifyInitCmd.Run()

	spicetifyBackupCreate := exec.Command("spicetify", "backup", "apply")
	LogCommand("spicetify backup apply")
	err = spicetifyBackupCreate.Run()
	if err != nil {
		LogError("Spicetify backup/apply failed, but continuing", err)
	}

	spicetifyConfig := exec.Command("spicetify", "config", "current_theme", "Sleek")
	LogCommand("spicetify config current_theme Sleek")
	err = spicetifyConfig.Run()
	if err != nil {
		LogError("Spicetify theme config failed, but continuing", err)
	}

	spicetifyApply := exec.Command("spicetify", "apply")
	LogCommand("spicetify apply")
	err = spicetifyApply.Run()
	if err != nil {
		LogError("Spicetify apply failed, but continuing", err)
	}

	LogInfo("Spotify ricing complete")
	fmt.Println("Spotify ricing complete")
}

func copyFromUsrShareToLocalAndPerformOverwrite(ogFileName, keyString, replaceValue string) {
	homedir, err := os.UserHomeDir()
	ReturnOnErr(err)

	localAppsDir := homedir + "/.local/share/applications/"
	err = os.MkdirAll(localAppsDir, 0o755)
	ReturnOnErr(err)

	dstPath := localAppsDir + ogFileName

	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		srcPath := "/usr/share/applications/" + ogFileName
		if _, err := os.Stat(srcPath); os.IsNotExist(err) {
			fmt.Printf("Warning: Source file %s does not exist, skipping...\n", srcPath)
			return
		}

		srcFile, err := os.Open(srcPath)
		ReturnOnErr(err)
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		ReturnOnErr(err)
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		ReturnOnErr(err)
		srcFile.Close()
		dstFile.Close()
	}

	content, err := os.ReadFile(dstPath)
	if err != nil {
		fmt.Printf("Warning: Could not read %s, skipping...\n", dstPath)
		return
	}

	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		var split []string
		split = strings.Split(line, "=")
		if len(split) >= 1 && strings.TrimSpace(split[0]) == keyString {
			lines[i] = replaceValue
		}
	}

	err = os.WriteFile(dstPath, []byte(strings.Join(lines, "\n")), 0o644)
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

	exec.Command("rm", "-rf", "/tmp/yay-install").Run()

	return true
}

func installSddmTheme() {
	fmt.Println("Installing SDDM Astronaut Theme...")

	themeDir := "/usr/share/sddm/themes/sddm-astronaut-theme"

	if _, err := os.Stat(themeDir); err == nil {
		fmt.Println("SDDM Astronaut Theme already exists, skipping clone...")
	} else {
		cloneTheme := exec.Command("sudo", "git", "clone", "-b", "master", "--depth", "1",
			"https://github.com/keyitdev/sddm-astronaut-theme.git",
			themeDir)
		err := cloneTheme.Run()
		ReturnOnErr(err)
	}

	copyFonts := exec.Command("sudo", "cp", "-r",
		"/usr/share/sddm/themes/sddm-astronaut-theme/Fonts/",
		"/usr/share/fonts/")
	err := copyFonts.Run()
	ReturnOnErr(err)

	sddmConfig := `[Theme]
Current=sddm-astronaut-theme`
	err = os.WriteFile("/tmp/sddm.conf", []byte(sddmConfig), 0o644)
	ReturnOnErr(err)

	copySddmConfig := exec.Command("sudo", "cp", "/tmp/sddm.conf", "/etc/sddm.conf")
	err = copySddmConfig.Run()
	ReturnOnErr(err)

	virtKbdConfig := `[General]
InputMethod=qtvirtualkeyboard`
	err = os.WriteFile("/tmp/virtualkbd.conf", []byte(virtKbdConfig), 0o644)
	ReturnOnErr(err)

	mkdirSddm := exec.Command("sudo", "mkdir", "-p", "/etc/sddm.conf.d")
	err = mkdirSddm.Run()
	ReturnOnErr(err)

	copyVirtKbd := exec.Command("sudo", "cp", "/tmp/virtualkbd.conf", "/etc/sddm.conf.d/virtualkbd.conf")
	err = copyVirtKbd.Run()
	ReturnOnErr(err)

	enableSddm := exec.Command("sudo", "systemctl", "enable", "sddm")
	err = enableSddm.Run()
	ReturnOnErr(err)

	os.Remove("/tmp/sddm.conf")
	os.Remove("/tmp/virtualkbd.conf")

	fmt.Println("SDDM Astronaut Theme has been installed and configured.")
}

func setupWallpapers(homedir string) {
	fmt.Println("Setting up wallpapers...")

	wd, err := os.Getwd()
	ReturnOnErr(err)

	assetsWallpaperPath := wd + "/assets/wallpaper"
	homeWallpaperPath := homedir + "/wallpaper"

	if _, err := os.Stat(assetsWallpaperPath); os.IsNotExist(err) {
		fmt.Printf("Warning: Assets wallpaper directory %s does not exist, skipping wallpaper setup\n", assetsWallpaperPath)
		return
	}

	err = os.MkdirAll(homeWallpaperPath, 0o755)
	ReturnOnErr(err)

	copyWallpapers := exec.Command("cp", "-r", assetsWallpaperPath+"/.", homeWallpaperPath+"/")
	err = copyWallpapers.Run()
	ReturnOnErr(err)

	fmt.Printf("Wallpapers merged from %s to %s\n", assetsWallpaperPath, homeWallpaperPath)
}

func setupBashrc(homedir string) {
	fmt.Println("Setting up .bashrc...")

	bashrcLoop := `
	for f in ~/.config/bashrc/*; do
		if [ ! -d $f ] ;then
			c=` + "`echo $f | sed -e \"s=.config/bashrc=.config/bashrc/custom=\"`" + `
				[[ -f $c ]] && source $c || source $f
			fi
	done`

	bashrcPath := homedir + "/.bashrc"

	if _, err := os.Stat(bashrcPath); os.IsNotExist(err) {
		err := os.WriteFile(bashrcPath, []byte(bashrcLoop), 0o644)
		ReturnOnErr(err)
		fmt.Printf(".bashrc created at %s\n", bashrcPath)
	} else {
		content, err := os.ReadFile(bashrcPath)
		ReturnOnErr(err)

		contentStr := string(content)

		if strings.Contains(contentStr, "~/.config/bashrc/*") {
			fmt.Println(".bashrc already contains nxtdots configuration, skipping...")
			return
		}

		updatedContent := contentStr + bashrcLoop
		err = os.WriteFile(bashrcPath, []byte(updatedContent), 0o644)
		ReturnOnErr(err)
		fmt.Printf("Added nxtdots configuration to existing .bashrc at %s\n", bashrcPath)
	}
}

func setupOneTimeCommands() {
	fmt.Println("Running one-time setup commands...")

	setColorScheme := exec.Command("gsettings", "set", "org.gnome.desktop.interface", "color-scheme", "prefer-dark")
	err := setColorScheme.Run()
	ReturnOnErr(err)

	enableNetworkManager := exec.Command("sudo", "systemctl", "enable", "NetworkManager")
	err = enableNetworkManager.Run()
	ReturnOnErr(err)

	hyprsunsetInstall := exec.Command("systemctl", "--user", "enable", "--now", "hyprsunset.service")
	err = hyprsunsetInstall.Run()
	ReturnOnErr(err)

	fmt.Println("One-time setup commands completed")
}

func setupDesktopEntries(homedir string) {
	fmt.Println("Setting up desktop entries...")

	// Create necessary directories
	iconsDir := homedir + "/.local/share/icons"
	err := os.MkdirAll(iconsDir, 0o755)
	ReturnOnErr(err)

	claudeChatsDir := homedir + "/.cache/claude-chats"
	err = os.MkdirAll(claudeChatsDir, 0o755)
	ReturnOnErr(err)

	applicationsDir := homedir + "/.local/share/applications"
	err = os.MkdirAll(applicationsDir, 0o755)
	ReturnOnErr(err)

	// Copy Claude.png to icons directory
	configDir := homedir + "/.config"
	claudeIconSrc := configDir + "/assets/icons/Claude.png"
	claudeIconDst := iconsDir + "/Claude.png"

	if _, err := os.Stat(claudeIconSrc); err == nil {
		srcFile, err := os.Open(claudeIconSrc)
		if err == nil {
			defer srcFile.Close()
			dstFile, err := os.Create(claudeIconDst)
			if err == nil {
				defer dstFile.Close()
				_, err = io.Copy(dstFile, srcFile)
				ReturnOnErr(err)
				fmt.Printf("Copied Claude icon to %s\n", claudeIconDst)
			}
		}
	} else {
		fmt.Printf("Warning: Claude.png not found at %s, skipping icon installation\n", claudeIconSrc)
	}

	// Create Claude Code desktop entry
	claudeDesktopEntry := `[Desktop Entry]
Name=Claude Code
Comment=AI-powered coding assistant by Anthropic
Exec=claude-code
Icon=` + iconsDir + `/Claude.png
Type=Application
Categories=Development;Utility;
Terminal=false
StartupNotify=true`

	claudeDesktopPath := applicationsDir + "/claude-code.desktop"
	err = os.WriteFile(claudeDesktopPath, []byte(claudeDesktopEntry), 0o644)
	ReturnOnErr(err)
	fmt.Printf("Created Claude Code desktop entry at %s\n", claudeDesktopPath)

	// Create Neovim desktop entry
	nvimDesktopEntry := `[Desktop Entry]
Name=Neovim
Comment=Hyperextensible Vim-based text editor
Exec=kitty nvim %F
Icon=nvim
Type=Application
Categories=Development;TextEditor;
Terminal=false
StartupNotify=true
MimeType=text/plain;text/x-makefile;text/x-c++hdr;text/x-c++src;text/x-chdr;text/x-csrc;text/x-java;text/x-moc;text/x-pascal;text/x-tcl;text/x-tex;application/x-shellscript;text/x-c;text/x-c++;`

	nvimDesktopPath := applicationsDir + "/nvim.desktop"
	err = os.WriteFile(nvimDesktopPath, []byte(nvimDesktopEntry), 0o644)
	ReturnOnErr(err)
	fmt.Printf("Created Neovim desktop entry at %s\n", nvimDesktopPath)

	pacseekDesktopEntry := `[Desktop Entry]
Name=Pacseek
Comment=Terminal UI for searching and installing Arch Linux packages
Exec=kitty --class Pacseek pacseek
Icon=system-software-install
Type=Application
Categories=System;PackageManager;
Terminal=false
StartupNotify=true`

	pacseekDesktopPath := applicationsDir + "/org.moson.pacseek.desktop"
	err = os.WriteFile(pacseekDesktopPath, []byte(pacseekDesktopEntry), 0o644)
	ReturnOnErr(err)
	fmt.Printf("Created Pacseek desktop entry at %s\n", pacseekDesktopPath)

	fmt.Println("Desktop entries setup complete")
}

func updateWaybarThemePaths(homedir string) {
	fmt.Println("Updating waybar theme paths...")

	darkThemePath := homedir + "/.config/waybar/themes/nxtdots-pywal-dark.css"
	updateCSSFilePath(darkThemePath, homedir)

	lightThemePath := homedir + "/.config/waybar/themes/nxtdots-pywal-light.css"
	updateCSSFilePath(lightThemePath, homedir)

	fmt.Println("Waybar theme paths updated")
}

func updateCSSFilePath(filePath, homedir string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Warning: Could not read %s, skipping path update\n", filePath)
		return
	}

	oldPath := "file:///home/nxtdots/.cache/wal/colors-waybar.css"
	newPath := "file://" + homedir + "/.cache/wal/colors-waybar.css"
	updatedContent := strings.ReplaceAll(string(content), oldPath, newPath)

	err = os.WriteFile(filePath, []byte(updatedContent), 0o644)
	ReturnOnErr(err)

	fmt.Printf("Updated path in %s\n", filePath)
}

func updateWaybarMainTheme(homedir string) {
	fmt.Println("Updating main waybar theme path...")

	mainThemePath := homedir + "/.config/waybar/theme.css"
	updateCSSFilePath(mainThemePath, homedir)

	fmt.Println("Main waybar theme path updated")
}

func installPackageGroup(packages []string) {
	if len(packages) == 0 {
		return
	}

	LogInfo(fmt.Sprintf("Installing %d packages: %s", len(packages), strings.Join(packages, ", ")))

	cmd := exec.Command("yay", append([]string{"-S", "--needed", "--noconfirm"}, packages...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		LogError(fmt.Sprintf("Package installation failed for group (some packages may have been installed): %v", packages), err)
	} else {
		LogInfo(fmt.Sprintf("Successfully completed installation attempt for %d packages", len(packages)))
	}
}

func isPackageInstalled(packageName string) bool {
	cmd := exec.Command("pacman", "-Qi", packageName)
	err := cmd.Run()
	return err == nil
}

func setupNvidiaIfDetected(homedir string) {
	fmt.Println("Checking for NVIDIA GPU...")

	cmd := exec.Command("lspci")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Warning: Could not detect GPU: %v\n", err)
		return
	}

	outputStr := strings.ToLower(string(output))
	if !strings.Contains(outputStr, "nvidia") {
		fmt.Println("No NVIDIA GPU detected, skipping NVIDIA setup")
		return
	}

	fmt.Println("NVIDIA GPU detected! Installing NVIDIA drivers...")
	LogInfo("NVIDIA GPU detected, installing drivers")

	nvidiaPackages := []string{
		"nvidia",
		"nvidia-utils",
		"egl-wayland",
		"libva-nvidia-driver",
	}

	installPackageGroup(nvidiaPackages)

	fmt.Println("Configuring GRUB for NVIDIA...")
	LogInfo("Configuring GRUB for NVIDIA")

	grubPath := "/etc/default/grub"
	grubContent, err := os.ReadFile(grubPath)
	if err != nil {
		LogError("Failed to read GRUB config", err)
		fmt.Printf("Warning: Could not read %s, please manually add nvidia_drm.modeset=1 to GRUB_CMDLINE_LINUX_DEFAULT\n", grubPath)
	} else {
		grubStr := string(grubContent)
		if !strings.Contains(grubStr, "nvidia_drm.modeset=1") {
			lines := strings.Split(grubStr, "\n")
			for i, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "GRUB_CMDLINE_LINUX_DEFAULT=") {
					line = strings.TrimRight(line, "\"")
					if !strings.HasSuffix(line, "=") && !strings.HasSuffix(line, "=\"") {
						line += " "
					}
					lines[i] = line + "nvidia_drm.modeset=1 nvidia_drm.fbdev=1\""
					break
				}
			}
			updatedGrub := strings.Join(lines, "\n")
			tmpGrubPath := "/tmp/grub-nvidia"
			err = os.WriteFile(tmpGrubPath, []byte(updatedGrub), 0o644)
			if err == nil {
				copyCmd := exec.Command("sudo", "cp", tmpGrubPath, grubPath)
				err = copyCmd.Run()
				if err != nil {
					LogError("Failed to update GRUB config", err)
				} else {
					fmt.Println("GRUB config updated with NVIDIA parameters")
					LogInfo("GRUB config updated successfully")

					grubMkconfig := exec.Command("sudo", "grub-mkconfig", "-o", "/boot/grub/grub.cfg")
					grubMkconfig.Stdout = os.Stdout
					grubMkconfig.Stderr = os.Stderr
					err = grubMkconfig.Run()
					if err != nil {
						LogError("Failed to regenerate GRUB config", err)
					} else {
						fmt.Println("GRUB configuration regenerated")
					}
				}
				os.Remove(tmpGrubPath)
			}
		} else {
			fmt.Println("GRUB already configured for NVIDIA")
		}
	}

	fmt.Println("Configuring Hyprland environment for NVIDIA...")
	LogInfo("Configuring Hyprland environment for NVIDIA")

	envConfPath := homedir + "/.config/hypr/conf/environment.conf"
	envContent, err := os.ReadFile(envConfPath)
	if err != nil {
		LogError("Failed to read environment.conf", err)
		fmt.Printf("Warning: Could not read %s\n", envConfPath)
		return
	}

	envStr := string(envContent)
	lines := strings.Split(envStr, "\n")
	modified := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# env = GBM_BACKEND,nvidia-drm") ||
			strings.HasPrefix(trimmed, "# env = LIBVA_DRIVER_NAME,nvidia") ||
			strings.HasPrefix(trimmed, "# env = __GLX_VENDOR_LIBRARY_NAME,nvidia") ||
			strings.HasPrefix(trimmed, "# env = __NV_PRIME_RENDER_OFFLOAD,1") ||
			strings.HasPrefix(trimmed, "# env = __VK_LAYER_NV_optimus,NVIDIA_only") {
			lines[i] = strings.Replace(line, "# env =", "env =", 1)
			modified = true
		}
		if strings.HasPrefix(trimmed, "# cursor {") {
			lines[i] = "cursor {"
			modified = true
		}
		if strings.HasPrefix(trimmed, "#     no_hardware_cursors = true") {
			lines[i] = "    no_hardware_cursors = true"
			modified = true
		}
		if strings.HasPrefix(trimmed, "# }") && modified {
			lines[i] = "}"
		}
	}

	if modified {
		updatedEnv := strings.Join(lines, "\n")
		err = os.WriteFile(envConfPath, []byte(updatedEnv), 0o644)
		if err != nil {
			LogError("Failed to update environment.conf", err)
		} else {
			fmt.Println("Hyprland environment configured for NVIDIA")
			LogInfo("environment.conf updated successfully")
		}
	} else {
		fmt.Println("Hyprland environment already configured for NVIDIA")
	}

	fmt.Println("NVIDIA setup complete! Please reboot your system for changes to take effect.")
	LogInfo("NVIDIA setup completed")
}

func runInitialPywalSetup(homedir string) {
	fmt.Println("Running initial pywal setup...")

	wallpaperDir := homedir + "/wallpaper"
	entries, err := os.ReadDir(wallpaperDir)
	if err != nil {
		LogError("Failed to read wallpaper directory", err)
		fmt.Printf("Warning: Could not read wallpaper directory %s\n", wallpaperDir)
		return
	}

	var firstWallpaper string
	for _, entry := range entries {
		if !entry.IsDir() {
			name := strings.ToLower(entry.Name())
			if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") ||
				strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".webp") {
				firstWallpaper = wallpaperDir + "/" + entry.Name()
				break
			}
		}
	}

	if firstWallpaper == "" {
		fmt.Println("No wallpaper found in ~/wallpaper, skipping initial pywal setup")
		LogInfo("No wallpaper found, skipping pywal setup")
		return
	}

	fmt.Printf("Generating initial theme from wallpaper: %s\n", firstWallpaper)
	LogInfo(fmt.Sprintf("Running initial pywal with wallpaper: %s", firstWallpaper))

	colorScheme := GetCurrentSystemTheme()
	err = ExecPywal(colorScheme, firstWallpaper)
	if err != nil {
		LogError("Failed to run pywal", err)
		fmt.Printf("Warning: pywal failed: %v\n", err)
		return
	}

	colors, err := GetPywalColors()
	if err != nil {
		LogError("Failed to get pywal colors", err)
		fmt.Printf("Warning: Could not get pywal colors: %v\n", err)
		return
	}

	cavaConfigPath := homedir + "/.config/cava/config"
	err = UpdateCavaGradient(cavaConfigPath, colors)
	if err != nil {
		LogError("Failed to update cava gradient", err)
		fmt.Printf("Warning: Could not update cava gradient: %v\n", err)
	} else {
		fmt.Println("Cava gradient configured")
	}

	cachePath := homedir + "/.cache/wallpaper/"
	err = os.MkdirAll(cachePath+"wallpaper-generated/", 0o755)
	if err == nil {
		os.WriteFile(cachePath+"current_wallpaper", []byte(firstWallpaper), 0o644)
	}

	fmt.Println("Initial pywal setup complete")
	LogInfo("Initial pywal setup completed successfully")
}

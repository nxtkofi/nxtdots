package utils

import (
	"os"
	"os/exec"
	"strings"
)

func HandleThemeChange(option, homeDir string) error {
	currentSystemTheme := GetCurrentSystemTheme()
	currentWallpaperFullPath := GetCurrentWallpaperFullPath(homeDir)

	switch option {
	case "toggle":
		if currentSystemTheme == PrefersDark {
			toggleThemeColor(PrefersLight, string(currentWallpaperFullPath), homeDir)
		} else {
			toggleThemeColor(PrefersDark, string(currentWallpaperFullPath), homeDir)
		}

	}

	return nil
}

func toggleThemeColor(prefersNewColorMode XDGColorScheme, currentWallpaperFullPath, homeDir string) {
	SetCurrentSystemTheme(prefersNewColorMode)
	err := ExecPywal(prefersNewColorMode, currentWallpaperFullPath)
	ReturnOnErr(err)
	colors, err := GetPywalColors()
	ReturnOnErr(err)
	UpdateSpicetify(colors, homeDir)
	updateGtkTheme(prefersNewColorMode)
	updateWalcord(homeDir)
	resetKittyIfItsRunning()

	themeDirPath := homeDir + "/.config/waybar/themes"
	themeDirCss := homeDir + "/.config/waybar/theme.css"
	files, err := os.ReadDir(themeDirPath)
	ReturnOnErr(err)
	var fileThemeModeToCopy string
	if prefersNewColorMode == PrefersDark {
		fileThemeModeToCopy = "dark"
	} else {
		fileThemeModeToCopy = "light"
	}

	for _, file := range files {
		if strings.Contains(file.Name(), fileThemeModeToCopy) {
			err = copyFile(themeDirPath+"/"+file.Name(), themeDirCss)
			ReturnOnErr(err)
		}
	}
}

func updateWalcord(homeDir string) {
	templatePath := homeDir + "/.config/vesktop/themes/midnight-vesktop.template.css"
	outputPath := homeDir + "/.config/vesktop/themes/midnight-vesktop.theme.css"
	cmd := exec.Command("walcord", "-t", templatePath, "-o", outputPath)
	err := cmd.Run()
	ReturnOnErr(err)
}

//	NOTE: note that below function seems counter-intuitive.
//
// I don't understand it, I have failed to find an answer as to why this is the only config that works, I have no clue and no longer any desire to find an answer
func updateGtkTheme(desiredColorScheme XDGColorScheme) {
	if desiredColorScheme == PrefersDark {
		cmd := exec.Command("gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", "'Adwaita'")
		err := cmd.Run()
		ReturnOnErr(err)
	} else {
		cmd := exec.Command("gsettings", "set", "org.gnome.desktop.interface", "gtk-theme", "'Adwaita-dark'")
		err := cmd.Run()
		ReturnOnErr(err)
	}
}

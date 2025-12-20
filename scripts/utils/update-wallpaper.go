package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

func UpdateWallpaper(newWallpaperFullFilePath, homeDir string) {
	homeDir, err := os.UserHomeDir()
	ReturnOnErr(err)
	cavaConfigPath := homeDir + "/.config/cava/config"
	cachePath := homeDir + "/.cache/wallpaper/"

	err = os.MkdirAll(filepath.Dir(cachePath+"wallpaper-generated/"), os.ModePerm)

	ReturnOnErr(err)
	err = os.WriteFile(cachePath+"current_wallpaper", []byte(newWallpaperFullFilePath), os.ModePerm)

	ReturnOnErr(err)

	colorScheme := GetCurrentSystemTheme()

	swwwCmd := exec.Command("swww", "img", newWallpaperFullFilePath,
		"--transition-type", "grow",
		"--transition-pos", "0.5,0.5",
		"--transition-duration", "1.5",
		"--transition-fps", "165",
		"--transition-bezier", "0.25,0.1,0.25,1.0")
	err = swwwCmd.Run()
	ReturnOnErr(err)

	err = ExecPywal(colorScheme, newWallpaperFullFilePath)
	colors, err := GetPywalColors()
	ReturnOnErr(err)

	err = UpdateCavaGradient(cavaConfigPath, colors)

	ReturnOnErr(err)
	err = UpdateSpicetify(colors, homeDir)
	ReturnOnErr(err)

	walcordUpdate := exec.Command("walcord", "-j", homeDir+"/.cache/wal/colors.json", "-t", homeDir+"/.config/vesktop/themes/midnight-vesktop.template.css", "-o", homeDir+"/.config/vesktop/themes/midnight-vesktop.theme.css")
	err = walcordUpdate.Run()

	ReturnOnErr(err)

	swayncUpdate := exec.Command("swaync-client", "-rs")
	err = swayncUpdate.Run()
	ReturnOnErr(err)

	err = GetOrCreateWallpaperCache(homeDir, newWallpaperFullFilePath)
	ReturnOnErr(err)

	err = RestartWaybar()
	ReturnOnErr(err)
}

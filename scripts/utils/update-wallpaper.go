package utils

import (
	"os"
	"os/exec"
	"path/filepath"
)

func UpdateWallpaper(newWallpaperFullFilePath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	cavaConfigPath := homeDir + "/.config/cava/config"
	cachePath := homeDir + "/.cache/wallpaper/"

	err = os.MkdirAll(filepath.Dir(cachePath+"wallpaper-generated/"), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(cachePath+"current_wallpaper", []byte(newWallpaperFullFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	err = ExecPywal(newWallpaperFullFilePath)
	colors, err := GetPywalColors()

	if err != nil {
		return err
	}
	err = UpdateCavaGradient(cavaConfigPath, colors)
	if err != nil {
		return err
	}

	err = UpdateSpicetify(colors, homeDir)
	if err != nil {
		return err
	}

	walcordUpdate := exec.Command("walcord", "-i", newWallpaperFullFilePath, "-t", homeDir+"/.config/vesktop/themes/midnight-vesktop.template.css", "-o", homeDir+"/.config/vesktop/themes/midnight-vesktop.theme.css")
	err = walcordUpdate.Run()
	if err != nil {
		return err
	}

	swayncUpdate := exec.Command("swaync-client", "-rs")
	err = swayncUpdate.Run()
	if err != nil {
		return err
	}
	err = GetOrCreateWallpaperCache(homeDir, newWallpaperFullFilePath)

	return nil
}

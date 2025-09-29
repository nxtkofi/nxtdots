package utils

import (
	"io"
	"os"
	"os/exec"
	"strings"
)

func GetOrCreateWallpaperCache(homeDir string, wallpaperFullPath string) error {

	splitWallpaperPath := strings.Split(wallpaperFullPath, "/")
	wallpaperFileName := splitWallpaperPath[len(splitWallpaperPath)-1]
	blurredWallpaperPath := homeDir + "/.cache/wallpaper/blurred_wallpaper.png"

	err := readOrCreateBlurredWallpaper(homeDir, wallpaperFileName, wallpaperFullPath, blurredWallpaperPath)
	if err != nil {
		return err
	}

	err = createRasiFile(homeDir+"/.cache/wallpaper/current_wallpaper.rasi", blurredWallpaperPath)
	if err != nil {
		return err
	}

	err = readOrCreateSquareWallpaper(homeDir, wallpaperFileName, wallpaperFullPath)
	if err != nil {
		return err
	}

	return nil

}
func readOrCreateBlurredWallpaper(homeDir, wallpaperFileName, wallPaperFullPath, destPath string) error {
	blur := "50x30"
	sourcePath := homeDir + "/.cache/wallpaper/wallpaper-generated/blur-" + blur + "-" + wallpaperFileName + ".png"

	if _, err := os.Stat(sourcePath); err != nil {
		cmd := exec.Command("magick", wallPaperFullPath, "-resize", "75%", destPath)
		err = cmd.Run()
		if err != nil {
			return err
		}
		if blur != "0x0" {
			cmd := exec.Command("magick", destPath, "-blur", blur, destPath)
			err = cmd.Run()
			if err != nil {
				return err
			}
		}
		err = copyFile(destPath, sourcePath)
		if err != nil {
			return err
		}
	}

	return copyFile(sourcePath, destPath)
}

func readOrCreateSquareWallpaper(homeDir, wallpaperFileName, wallpaperFullPath string) error {
	sourcePath := homeDir + "/.cache/wallpaper/wallpaper-generated/square-" + wallpaperFileName + ".png"
	destPath := homeDir + "/.cache/wallpaper/square_wallpaper.png"

	if _, err := os.Stat(sourcePath); err != nil {
		cmd := exec.Command("magick", wallpaperFullPath, "-gravity", "Center", "-extent", "1:1", destPath)
		err = cmd.Run()
		if err != nil {
			return err
		}

		err = copyFile(destPath, sourcePath)
		if err != nil {
			return err
		}
	}

	return copyFile(sourcePath, destPath)

}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

func createRasiFile(dst, blurredWallpaperPath string) error {
	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, err = file.WriteString("* { current-image: url(\"" + blurredWallpaperPath + "\", height); }")
	if err != nil {
		return err
	}
	return nil

}

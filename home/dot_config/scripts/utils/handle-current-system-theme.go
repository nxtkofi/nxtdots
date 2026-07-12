package utils

import (
	"os/exec"
	"strings"
)

type XDGColorScheme string

const (
	PrefersDark  XDGColorScheme = "'prefer-dark'"
	PrefersLight XDGColorScheme = "'prefer-light'"
)

func GetCurrentSystemTheme() XDGColorScheme {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	output, err := cmd.Output()
	cleanOutput := strings.TrimSpace(string(output))
	ReturnOnErr(err)
	return XDGColorScheme(cleanOutput)
}

func SetCurrentSystemTheme(pref XDGColorScheme) {
	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.interface", "color-scheme", string(pref))
	err := cmd.Run()
	ReturnOnErr(err)
}

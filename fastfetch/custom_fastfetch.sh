#!/bin/bash
THEME=$(gsettings get org.gnome.desktop.interface color-scheme)

LOGO_LIGHT="$HOME/.config/fastfetch/logos/sentry.png"
LOGO_DARK="$HOME/.config/fastfetch/logos/priestess.png"

LOGO_PATH=$LOGO_LIGHT

if [[ "$THEME" == *'dark'* ]]; then
  LOGO_PATH=$LOGO_DARK
else
  LOGO_PATH=$LOGO_LIGHT
fi

fastfetch -c "$HOME/.config/fastfetch/config.jsonc" -l "$LOGO_PATH"

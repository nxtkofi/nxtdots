#!/usr/bin/env bash

theme_css="$HOME/.config/waybar/theme.css"
theme_dir="$HOME/.config/waybar/themes"

current_theme=$(head -n 1 "$theme_css" | awk '{print $2}')


tooltip=$(tr '-' ' ' <<<"$current_theme")
echo "{ \"text\": \">\", \"tooltip\": \"Theme: <span text_transform='capitalize'>${tooltip}</span>\" }"

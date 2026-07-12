#!/bin/bash
cliphist list | rofi -dmenu -replace -config ~/.config/rofi/config-cliphist.rasi | cliphist decode | wl-copy

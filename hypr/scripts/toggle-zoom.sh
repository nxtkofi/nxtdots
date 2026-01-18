#!/bin/bash
cache_file="$HOME/.cache/toggle_zoom"
if [ -f $cache_file ]; then
    hyprctl keyword cursor:zoom_factor 1.0
    rm $cache_file
else
    hyprctl keyword cursor:zoom_factor 3.0
    touch $cache_file
fi
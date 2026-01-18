#!/bin/bash
defaultwallpaper="$HOME/wallpaper/default.jpg"
cachefile="$HOME/.cache/wallpaper/current_wallpaper"
if [ -f "$cachefile" ]; then
    sed -i "s|~|$HOME|g" "$cachefile"
    wallpaper=$(cat $cachefile)
    if [ -f $wallpaper ]; then
        echo ":: Wallpaper $wallpaper exists"
    else
        wallpaper=$defaultwallpaper
    fi
else
    wallpaper=$defaultwallpaper
fi

swww img "$wallpaper" --transition-type grow --transition-pos 0.5,0.5 --transition-duration 1.5 --transition-fps 165 --transition-bezier 0.25,0.1,0.25,1.0

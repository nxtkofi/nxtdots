#!/bin/bash
defaultwallpaper="$HOME/wallpaper/default.jpg"
cachefile="$HOME/.config/cache/current_wallpaper"
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

waypaper --wallpaper "$wallpaper"

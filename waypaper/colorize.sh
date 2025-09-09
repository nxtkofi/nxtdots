#!/bin/bash
if [ -f ~/.cache/wallpaper/cached_wallpaper ]; then
    use_cache=1
    echo ":: Using Wallpaper Cache"
else
    use_cache=0
    echo ":: Wallpaper Cache disabled"
fi

# -----------------------------------------------------
# Set defaults
# -----------------------------------------------------

force_generate=0
generatedversions="$HOME/.cache/wallpaper/wallpaper-generated"
waypaperrunning=$HOME/.cache/wallpaper/waypaper-running
cachefile="$HOME/.cache/wallpaper/current_wallpaper"
blurredwallpaper="$HOME/.cache/wallpaper/blurred_wallpaper.png"
squarewallpaper="$HOME/.cache/wallpaper/square_wallpaper.png"
rasifile="$HOME/.cache/wallpaper/current_wallpaper.rasi"
defaultwallpaper="$HOME/wallpaper/default.jpg"
wallpapereffect="off"
blur="50x30"

# Ensures that the script only run once if wallpaper effect enabled
if [ -f $waypaperrunning ]; then
    rm $waypaperrunning
    exit
fi

# Create folder with generated versions of wallpaper if not exists
if [ ! -d $generatedversions ]; then
    mkdir $generatedversions
fi

# -----------------------------------------------------
# Get selected wallpaper
# -----------------------------------------------------

if [ -z $1 ]; then
    if [ -f $cachefile ]; then
        wallpaper=$(cat $cachefile)
    else
        wallpaper=$defaultwallpaper
    fi
else
    wallpaper=$1
fi
used_wallpaper=$wallpaper
echo ":: Setting wallpaper with source image $wallpaper"
tmpwallpaper=$wallpaper

# -----------------------------------------------------
# Copy path of current wallpaper to cache file
# -----------------------------------------------------

if [ ! -f $cachefile ]; then
    touch $cachefile
fi
echo "$wallpaper" >$cachefile
echo ":: Path of current wallpaper copied to $cachefile"

# -----------------------------------------------------
# Get wallpaper filename
# -----------------------------------------------------
wallpaperfilename=$(basename $wallpaper)
echo ":: Wallpaper Filename: $wallpaperfilename"

# -----------------------------------------------------
# Execute pywal
# -----------------------------------------------------
color_scheme=$(gsettings get org.gnome.desktop.interface color-scheme)

if [[ "$color_scheme" == *'light'* ]]; then
  wal -q -l -i "$used_wallpaper" >> ~/.config/settings/post_debug.log 2>&1
else
  wal -q -i "$used_wallpaper" >> ~/.config/settings/post_debug.log 2>&1
fi

# -----------------------------------------------------
# Update Cava colors
# -----------------------------------------------------
echo ":: Updating Cava gradient"
$HOME/.config/settings/update-cava-gradient.sh

# -----------------------------------------------------
# Update Spicetify theme //TODO: Create spicetify config 
# pywal-spicetify text 
walcord -i "$wallpaper" -t ~/.config/vesktop/themes/midnight-vesktop.template.css -o ~/.config/vesktop/themes/midnight-vesktop.theme.css 

# -----------------------------------------------------
# Update Pywalfox
# -----------------------------------------------------

if type pywalfox >/dev/null 2>&1; then
    pywalfox update
fi

# -----------------------------------------------------
# Update SwayNC
# -----------------------------------------------------
sleep 0.1
swaync-client -rs

# -----------------------------------------------------
# Created blurred wallpaper
# -----------------------------------------------------

if [ -f $generatedversions/blur-$blur-$effect-$wallpaperfilename.png ] && [ "$force_generate" == "0" ] && [ "$use_cache" == "1" ]; then
    echo ":: Use cached wallpaper blur-$blur-$effect-$wallpaperfilename"
else
    echo ":: Generate new cached wallpaper blur-$blur-$effect-$wallpaperfilename with blur $blur"
    # notify-send --replace-id=1 "Generate new blurred version" "with blur $blur" -h int:value:66
    magick $used_wallpaper -resize 75% $blurredwallpaper
    echo ":: Resized to 75%"
    if [ ! "$blur" == "0x0" ]; then
        magick $blurredwallpaper -blur $blur $blurredwallpaper
        cp $blurredwallpaper $generatedversions/blur-$blur-$effect-$wallpaperfilename.png
        echo ":: Blurred"
    fi
fi
cp $generatedversions/blur-$blur-$effect-$wallpaperfilename.png $blurredwallpaper

# -----------------------------------------------------
# Create rasi file
# -----------------------------------------------------

if [ ! -f $rasifile ]; then
    touch $rasifile
fi
echo "* { current-image: url(\"$blurredwallpaper\", height); }" >"$rasifile"

# -----------------------------------------------------
# Created square wallpaper
# -----------------------------------------------------

echo ":: Generate new cached wallpaper square-$wallpaperfilename"
magick $tmpwallpaper -gravity Center -extent 1:1 $squarewallpaper
cp $squarewallpaper $generatedversions/square-$wallpaperfilename.png

# Send SIGUSR1 to all kitty processes to reload config
if pgrep -x "kitty" > /dev/null; then
    echo ":: Sending SIGUSR1 to all kitty processes"
    pkill -USR1 -x kitty
fi

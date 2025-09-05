# Fastfetch
if [[ $(tty) == *"pts"* ]]; then
    fastfetch --config examples/28
else
    echo
    if [ -f /bin/hyprctl ]; then
        echo "Start Hyprland with command Hyprland"
    fi
fi

# Starship
eval "$(starship init bash)"

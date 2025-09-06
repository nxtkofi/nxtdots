# Fastfetch
if [[ $(tty) == *"pts"* ]]; then
        $HOME/.config/fastfetch/custom_fastfetch.sh
else
    echo
    if [ -f /bin/hyprctl ]; then
        echo "Start Hyprland with command Hyprland"
    fi
fi

# Starship
eval "$(starship init bash)"

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

bashrc_sidecar="${BASH_SOURCE[0]%.sh}.local.sh"
[ -r "$bashrc_sidecar" ] && source "$bashrc_sidecar"
unset bashrc_sidecar

#!/bin/bash
sleep 2

# Restart the relevant services via systemd to ensure they are running correctly.
systemctl --user restart xdg-desktop-portal-hyprland
systemctl --user restart xdg-desktop-portal

# Update hyprsunset config with dynamic sunset time
SUNSET_TIME=$($HOME/.config/hypr/scripts/sunset.sh)
if [ -n "$SUNSET_TIME" ]; then
    sed -i "10s|time = .*|time = $SUNSET_TIME|" /home/nxtkofi/.config/hypr/hyprsunset.conf
fi

# Start hyprsunset
hyprsunset &

# Run waybar
waybar &

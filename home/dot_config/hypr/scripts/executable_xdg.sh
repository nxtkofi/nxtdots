#!/bin/bash
sleep 2

# Restart the relevant services via systemd to ensure they are running correctly.
systemctl --user restart xdg-desktop-portal-hyprland
systemctl --user restart xdg-desktop-portal

SUNSET_SCRIPT="$HOME/.config/hypr/scripts/sunset.sh"
SUNSET_TIME="17:00"
if [ -x "$SUNSET_SCRIPT" ]; then
    GENERATED_SUNSET_TIME=$("$SUNSET_SCRIPT")
    if [ -n "$GENERATED_SUNSET_TIME" ]; then
        SUNSET_TIME="$GENERATED_SUNSET_TIME"
    fi
fi

cat > "$HOME/.config/hypr/hyprsunset.conf" <<EOF
profile {
    time = 07:00
    identity = true
}

profile {
    time = $SUNSET_TIME
    temperature = 3500
}
EOF

# Start hyprsunset
hyprsunset &

# Run waybar
waybar &

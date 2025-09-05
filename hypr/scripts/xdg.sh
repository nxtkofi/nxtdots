#!/bin/bash
# __  ______   ____
# \ \/ /  _ \ / ___|
#  \  /| | | | |  _
#  /  \| |_| | |_| |
# /_/\_\____/ \____|
#

# Setup Timers
_sleep="2"

# Restart the relevant services via systemd to ensure they are running correctly.
systemctl --user restart xdg-desktop-portal-hyprland
systemctl --user restart xdg-desktop-portal

# Run waybar
waybar &

#!/bin/bash
current=$(hyprshade current)

if [ -z "$current" ]; then
    hyprshade on blue-light-filter-50
else
    hyprshade off
fi

#!/bin/bash

# Universal temperature script for waybar
# Works on both AMD (k10temp) and Intel (coretemp) systems

get_temp() {
    # Try k10temp (AMD)
    for hwmon in /sys/class/hwmon/hwmon*/name; do
        if [[ "$(cat "$hwmon")" == "k10temp" ]]; then
            hwmon_dir=$(dirname "$hwmon")
            if [[ -r "$hwmon_dir/temp1_input" ]]; then
                temp=$(cat "$hwmon_dir/temp1_input")
                temp_c=$((temp / 1000))
                echo "$temp_c"
                return 0
            fi
        fi
    done

    # Try coretemp (Intel)
    for hwmon in /sys/class/hwmon/hwmon*/name; do
        if [[ "$(cat "$hwmon")" == "coretemp" ]]; then
            hwmon_dir=$(dirname "$hwmon")
            if [[ -r "$hwmon_dir/temp1_input" ]]; then
                temp=$(cat "$hwmon_dir/temp1_input")
                temp_c=$((temp / 1000))
                echo "$temp_c"
                return 0
            fi
        fi
    done

    # Fallback to thermal_zone
    for zone in /sys/class/thermal/thermal_zone*/temp; do
        if [[ -r "$zone" ]]; then
            temp=$(cat "$zone")
            temp_c=$((temp / 1000))
            echo "$temp_c"
            return 0
        fi
    done

    echo "N/A"
}

temp=$(get_temp)
echo "{\"text\": \"$temp°C\", \"tooltip\": \"CPU Temperature: $temp°C\", \"class\": \"temperature\"}"
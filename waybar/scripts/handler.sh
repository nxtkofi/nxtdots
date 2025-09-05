#!/usr/bin/env bash

theme_css="$HOME/.config/waybar/theme.css"
theme_dir="$HOME/.config/waybar/themes"

current_theme=$(head -n 1 "$theme_css" | awk '{print $2}')

case $1 in
    'next' | 'prev')
        themes=("$theme_dir/"*.css)
        index=-1

        for i in "${!themes[@]}"; do
            theme=$(basename "${themes[$i]}" .css)
            if [[ $theme == "$current_theme" ]]; then
                index=$i
                break
            fi
        done

        if [[ $index -eq -1 ]]; then
            echo ":: Error: Could not find current theme '$current_theme' in theme list."
            exit 1
        fi

        case $1 in
            'next')
                new_index=$(((index + 1) % ${#themes[@]}))
                ;;
            'prev')
                new_index=$(((index - 1 + ${#themes[@]}) % ${#themes[@]}))
                ;;
        esac

        new_theme="${themes[$new_index]}"
        new_theme_name=$(basename "$new_theme" .css)

        if [[ "$new_theme_name" == *light* ]]; then
            gsettings set org.gnome.desktop.interface color-scheme 'prefer-light'
            
            wallpaper_path_raw=$(tail -n 2 "$HOME/.config/waypaper/config.ini" | head -n 1 | sed 's/wallpaper = //')
            wallpaper_path_expanded=$(eval echo "$wallpaper_path_raw")
            wal -q -i "$wallpaper_path_expanded" -l >> "$HOME/.config/settings/post_debug.log" 2>&1

            echo ":: Setting color-scheme to light"
        else
            gsettings set org.gnome.desktop.interface color-scheme 'prefer-dark'
            wallpaper_path_raw=$(tail -n 2 "$HOME/.config/waypaper/config.ini" | head -n 1 | sed 's/wallpaper = //')
            wallpaper_path_expanded=$(eval echo "$wallpaper_path_raw")
            wal -q -i "$wallpaper_path_expanded" >> "$HOME/.config/settings/post_debug.log" 2>&1
            echo ":: Setting color-scheme to dark"
        fi

        #zanim to sie wykona to wal musi juz wykonac swoja akcje.
        cp "$new_theme" "$theme_css"
        echo ":: Switching theme to $(basename "$new_theme" .css)"

        pkill waybar 2>/dev/null || true
        nohup waybar >/dev/null 2>&1 &
        ;;

    'toggle')
        killall -SIGUSR1 waybar
        ;;

    'reload')
        if pgrep -x waybar >/dev/null; then
            echo ":: Reloading Waybar..."
            killall -SIGUSR2 waybar
        else
            echo ":: Waybar is not running, nothing to reload"
        fi
        ;;
esac

tooltip=$(tr '-' ' ' <<<"$current_theme")
echo "{ \"text\": \">\", \"tooltip\": \"Theme: <span text_transform='capitalize'>${tooltip}</span>\" }"

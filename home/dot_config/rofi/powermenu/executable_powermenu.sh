#!/usr/bin/env bash

dir="$HOME/.config/rofi/powermenu"
theme='style'

# CMDs
current_time=$(date +"%H:%M")
host=$(hostname)

# Options
shutdown='󰐥'
reboot=''
lock=''
suspend='󰒲'
logout='󰍃'
yes='󰄬'
no='󰅖'

rofi_cmd() {
  rofi -dmenu \
    -p "Session Manager" \
    -mesg "${current_time} • Uptime: $(uptime -p | sed -e 's/up //g')" \
    -theme ${dir}/${theme}.rasi
}

confirm_cmd() {
  rofi -dmenu \
    -p 'Confirmation' \
    -mesg 'Are you Sure?' \
    -theme ${dir}/confirm.rasi
}

confirm_exit() {
  echo -e "$yes\n$no" | confirm_cmd
}

run_rofi() {
  echo -e "$lock\n$suspend\n$logout\n$reboot\n$shutdown" | rofi_cmd
}

run_cmd() {
  selected="$(confirm_exit)"
  if [[ "$selected" == "$yes" ]]; then
    if [[ $1 == '--shutdown' ]]; then
      systemctl poweroff
    elif [[ $1 == '--reboot' ]]; then
      systemctl reboot
    elif [[ $1 == '--suspend' ]]; then
      # Pause media playback
      playerctl pause 2>/dev/null
      # Mute audio
      wpctl set-mute @DEFAULT_AUDIO_SINK@ 1 2>/dev/null
      systemctl suspend
    elif [[ $1 == '--logout' ]]; then
      if [[ "$XDG_CURRENT_DESKTOP" == "Hyprland" ]]; then
        hyprctl dispatch exit
      elif [[ "$DESKTOP_SESSION" == 'openbox' ]]; then
        openbox --exit
      elif [[ "$DESKTOP_SESSION" == 'bspwm' ]]; then
        bspc quit
      elif [[ "$DESKTOP_SESSION" == 'i3' ]]; then
        i3-msg exit
      elif [[ "$DESKTOP_SESSION" == 'plasma' ]]; then
        qdbus org.kde.ksmserver /KSMServer logout 0 0 0
      fi
    fi
  else
    exit 0
  fi
}

chosen="$(run_rofi)"
case ${chosen} in
$shutdown)
  run_cmd --shutdown
  ;;
$reboot)
  run_cmd --reboot
  ;;
$lock)
  hyprlock
  ;;
$suspend)
  run_cmd --suspend
  ;;
$logout)
  run_cmd --logout
  ;;
esac

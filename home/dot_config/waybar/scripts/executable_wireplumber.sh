#!/usr/bin/env bash
#
# Control audio and microphone using wpctl
#
# Author: Jesse Mirabel <github.com/sejjy>
# Created: August 24, 2025
# License: MIT

LOCKFILE="/tmp/wireplumber_notify.lock"
DEBOUNCE_TIME=0.5

get_volume() {
	output=$(wpctl get-volume "$default")
	volume=$(awk '{print ($2 * 100)}' <<<"$output")
}

get_icon() {
	if ((volume == 0)); then
		icon="${prefix}-muted"
	elif ((volume < 33)); then
		icon="${prefix}-low"
	elif ((volume < 66)); then
		icon="${prefix}-medium"
	else
		icon="${prefix}-high"
	fi
}

notify_send() {
	get_volume "$1"
	get_icon "$1"

	case $output in
		*'MUTED'*)
			icon="${prefix}-muted"
			body='Muted'
			;;
		*)
			body='Unmuted'
			;;
	esac

	if [[ $2 == 'lower' ]] || [[ $2 == 'raise' ]]; then
		body="${volume}%"

		if [[ -f "$LOCKFILE" ]]; then
			kill "$(cat "$LOCKFILE")" 2>/dev/null
		fi

		(
			echo $$ > "$LOCKFILE"
			sleep "$DEBOUNCE_TIME"
			if [[ -f "$LOCKFILE" ]] && [[ "$(cat "$LOCKFILE")" == "$$" ]]; then
				notify-send "${title}: $body" -h int:value:"$volume" -i "$icon" -r 2425
				rm -f "$LOCKFILE"
			fi
		) &
	else
		notify-send "${title}: $body" -h int:value:"$volume" -i "$icon" -r 2425
	fi
}

case $1 in
	'mic')
		default='@DEFAULT_SOURCE@'
		prefix='microphone-sensitivity'
		title='Microphone'
		;;
	'out')
		default='@DEFAULT_SINK@'
		prefix='volume-level'
		title='Volume'
		;;
esac

value=5
case $2 in
	'lower') wpctl set-volume --limit 1.0 "$default" "$value"%- ;;
	'raise') wpctl set-volume --limit 1.0 "$default" "$value"%+ ;;
	*) wpctl set-mute "$default" toggle ;;
esac

notify_send "$1" "$2"

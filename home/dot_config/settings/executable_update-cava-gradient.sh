#!/bin/bash

colors_file="$HOME/.cache/wal/colors"
if [ ! -r "$colors_file" ]; then
  printf 'failed to read pywal colors\n' >&2
  exit 1
fi

colors=($(head -n 4 "$colors_file"))

theme_dir="$HOME/.config/cava/themes"
theme="$theme_dir/pywal.generated.local"
if ! mkdir -p "$theme_dir"; then
  printf 'failed to create Cava theme directory\n' >&2
  exit 1
fi

if ! cat > "$theme" <<EOF
[color]
gradient = 1
gradient_count = 3
gradient_color_1 = '${colors[1]}'
gradient_color_2 = '${colors[2]}'
gradient_color_3 = '${colors[3]}'
EOF
then
  printf 'failed to write Cava theme\n' >&2
  exit 1
fi

pkill -USR2 -x cava
status=$?
if [ "$status" -ne 0 ] && [ "$status" -ne 1 ]; then
  printf 'failed to reload Cava colors\n' >&2
  exit "$status"
fi

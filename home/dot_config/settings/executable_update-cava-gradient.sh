#!/bin/bash

colors=($(head -n 4 ~/.cache/wal/colors))

theme_dir="$HOME/.config/cava/themes"
theme="$theme_dir/pywal.generated.local"
mkdir -p "$theme_dir"
cat > "$theme" <<EOF
gradient = 1
gradient_count = 3
gradient_color_1 = '${colors[1]}'
gradient_color_2 = '${colors[2]}'
gradient_color_3 = '${colors[3]}'
EOF

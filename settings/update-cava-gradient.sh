#!/bin/bash

colors=($(head -n 4 ~/.cache/wal/colors))

cfg="$HOME/.config/cava/config"
sed -i "/^gradient_color_1/c\gradient_color_1 = '${colors[1]}'" "$cfg"
sed -i "/^gradient_color_2/c\gradient_color_2 = '${colors[2]}'" "$cfg"
sed -i "/^gradient_color_3/c\gradient_color_3 = '${colors[3]}'" "$cfg"

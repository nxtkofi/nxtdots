## Quickstart
```sh
git clone https://github.com/nxtkofi/nxtdots ~/.config/
cd .config
chmod +x ./install.sh
mkdir -p ~/.cache/wallpaper
./install.sh
```

## Introduction

Config was initially based on ml4w's dotfiles, grabbed a lot of basic stuff from
him, however now You probably would've never guessed that it was based on it.
Waybar is based on Mechabar, I added my own pywal-based-themes and modified them
(got rid of unused vars, basically same cleanup as I did in ml4w's case). I also
modified waybar's style a little bit and added custom theme switching script.

install.sh gets the necessary packages installed on Your system. Those are:
- fzf
- waybar 0.13.0 (0.14.0 has an issue rendering some things)
-

For floating pacseek (windowrule is already configured in
~/.config/hypr/conf/windowrule.conf) You have to edit org.moson.pacseek.desktop
file that's located in /usr/share/applications directory like so:
```.config
[Desktop Entry]

Name=pacseek
Comment=A terminal user interface for searching and installing Arch Linux packages

Icon=pacseek
Type=Application
Categories=Utility;
Keywords=terminal;package;

Exec=kitty --class Pacseek pacseek
StartupNotify=false
Terminal=false
```



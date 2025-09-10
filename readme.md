## Quickstart
```sh
git clone https://github.com/nxtkofi/nxtdots ~/.config/
cd .config
chmod +x ./install.sh
mkdir -p ~/.cache/wallpaper
./install.sh
```

## Requirements

Packages:
- fzf
- waybar 0.13.0 (0.14.0 has an issue rendering some things)
- vesktop (vencord, walcord)
- spicetify-cli
- python-pywal16

#### Commands:
Spotify rice:
- sudo chmod a+wr /opt/spotify
- sudo chmod a+wr /opt/spotify/Apps -R
one time:
`spicetify config current_theme Sleek`
`spicetify apply`


## Introduction

Config was initially based on ml4w's dotfiles, grabbed a lot of basic stuff from
him, however now You probably would've never guessed that it was based on it.
Waybar is based on Mechabar, I added my own pywal-based themes and modified them
(got rid of unused vars, basically same cleanup as I did in ml4w's case). I also
modified waybar's style a little bit and added custom theme switching script.

## Edited .desktop files

For floating pacseek (windowrule is already configured in
~/.config/hypr/conf/windowrule.conf) You have to edit org.moson.pacseek.desktop
file that's located in /usr/share/applications directory like so:
- pacseek
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

- vesktop
```
[Desktop Entry]
Name=Vesktop
Exec=/usr/bin/vesktop %U
Terminal=false
Type=Application
Icon=/home/<your-profile>/.config/assets/discord_custom.png
StartupWMClass=vesktop
GenericName=Internet Messenger
Categories=Network;
Keywords=discord;vencord;electron;chat;
Comment=Vesktop is a custom Discord App aiming to give you better performance and improve linux support. Vencord comes pre-installed
MimeType=x-scheme-handler/discord
```

- spotify
```
[Desktop Entry]
Type=Application
Name=Spotify
GenericName=Music Player
Icon=spotify-client
TryExec=spotify
Exec=/home/<your-profile>/.config/settings/launch-spotify.sh
Terminal=false
MimeType=x-scheme-handler/spotify;
Categories=Audio;Music;Player;AudioVideo;
StartupWMClass=spotify          
```

## Thanks
Thanks to ml4w for providing such a great base for this ricing
Big thanks to - https://www.reddit.com/r/unixporn/comments/1chv3tr/hyprland_everything_pywal/ (repo:https://github.com/magnusKue/wal-switcher/tree/master), he made pywal spicetifying easy!

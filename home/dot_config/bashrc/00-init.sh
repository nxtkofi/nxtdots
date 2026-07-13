export EDITOR=nvim
export PATH="/usr/lib/ccache/bin/:$PATH"
export PATH="$HOME/.local/bin:$PATH"
export PATH="$HOME/go/bin:$PATH"
export PATH="$HOME/.local/share/pnpm/bin:$PATH"
export NVM_DIR="$HOME/.nvm"
source /usr/share/nvm/init-nvm.sh

bashrc_sidecar="${BASH_SOURCE[0]%.sh}.local.sh"
[ -r "$bashrc_sidecar" ] && source "$bashrc_sidecar"
unset bashrc_sidecar

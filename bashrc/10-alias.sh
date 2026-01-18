alias c='clear'
alias wifi='nmtui'
alias vim='$EDITOR'

alias shutdown='systemctl poweroff'

if [[ $- == *i* ]]; then
    alias ls='eza --icons --color=always --group-directories-first'
else
    alias ls='ls --color=auto'
fi
alias grep='grep --color=auto'
alias cd='z'

alias mux='tmuxinator'

muxi() {
    local selected=$(tmuxinator list | tail -n +2 | xargs -n1 | fzf)
    if [ -n "$selected" ]; then
        tmuxinator start "$selected"
    fi
}

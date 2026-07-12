# chezmoi workflow

This repo is now managed by chezmoi.

## 1) Change only on this machine

1. Edit files directly in `~/.config/...`.
2. Sync local edits back to source:
   ```bash
   chezmoi re-add
   ```
3. Apply and verify:
   ```bash
   chezmoi diff
   ```

## 2) Change for your private NVIDIA laptop too

1. Edit files on one machine in `~/.config/...`.
2. Pull into source:
   ```bash
   chezmoi re-add
   ```
3. Commit and push from `~/.local/share/chezmoi`.
4. On the NVIDIA laptop, pull latest source and apply:
   ```bash
   chezmoi update
   ```
5. Keep laptop-specific values in `~/.config/chezmoi/chezmoi.toml`:
   - `has_nvidia = true`
   - laptop monitor lines in `monitors_raw`
   - laptop coordinates and secrets

## 3) Push changes to all public repo users

1. Edit shared files in `~/.config/...`.
2. Pull local changes into source:
   ```bash
   chezmoi re-add
   ```
3. Review and commit in `~/.local/share/chezmoi`.
4. Push to GitHub.
5. Other users run:
   ```bash
   chezmoi update
   ```

## Useful commands

```bash
chezmoi status
chezmoi diff
chezmoi doctor
chezmoi cd
```

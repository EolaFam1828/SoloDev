#!/usr/bin/env sh

set -eu

git config --local core.hooksPath .githooks
git config --local pull.ff only
git config --local pull.rebase false
git config --local fetch.prune true
git config --local merge.ff only
git config --local push.default simple
git config --local push.autoSetupRemote true
git config --local branch.autosetuprebase never
git config --local alias.sync '!git fetch origin --prune && git pull --ff-only origin main'

# Clear the stale template path from older local setup.
git config --local --unset-all commit.template 2>/dev/null || true

# Keep main explicitly tracking origin/main in fresh or repaired clones.
git branch --set-upstream-to=origin/main main >/dev/null 2>&1 || true

printf '%s\n' "Configured local Git defaults:"
printf '%s\n' "  hooksPath=.githooks"
printf '%s\n' "  pull.ff=only"
printf '%s\n' "  pull.rebase=false"
printf '%s\n' "  merge.ff=only"
printf '%s\n' "  fetch.prune=true"
printf '%s\n' "  push.default=simple"
printf '%s\n' "  push.autoSetupRemote=true"
printf '%s\n' "  alias.sync=git fetch origin --prune && git pull --ff-only origin main"

#!/bin/sh
# shellcheck disable=SC2148

# NOTE: If .env file doesn't exist, create a template file.
[[ -f .env ]] || tee .env >/dev/null <<'EOF'
# NOTE: Define environment variables that are referenced in the container.
EOF

# NOTE: Load .env files.
dotenv .versenv.env
dotenv .default.env
dotenv .env

# NOTE: Define environment variables that are NOT referenced in the container.
REPO_ROOT=$(git rev-parse --show-toplevel)
export REPO_ROOT
export PATH="${REPO_ROOT}/.local/bin:${REPO_ROOT}/.bin:${PATH}"
export DOCKER_BUILDKIT="1"
export COMPOSE_DOCKER_CLI_BUILD="1"

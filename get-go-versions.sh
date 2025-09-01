#!/usr/bin/env bash
set -euo pipefail

# ---------------------------------
# Config: repositories (GitHub URLs)
# ---------------------------------
REPOS=(
  "https://github.com/Netcracker/qubership-core-lib-go"
  "https://github.com/Netcracker/qubership-core-lib-go-actuator-common"
  "https://github.com/Netcracker/qubership-core-lib-go-bg-kafka"
  "https://github.com/Netcracker/qubership-core-lib-go-bg-state-monitor"
  "https://github.com/Netcracker/qubership-core-lib-go-dbaas-arangodb-client"
  "https://github.com/Netcracker/qubership-core-lib-go-dbaas-base-client"
  "https://github.com/Netcracker/qubership-core-lib-go-dbaas-cassandra-client"
  "https://github.com/Netcracker/qubership-core-lib-go-dbaas-clickhouse-client"
  "https://github.com/Netcracker/qubership-core-lib-go-dbaas-mongo-client"
  "https://github.com/Netcracker/qubership-core-lib-go-dbaas-opensearch-client"
  "https://github.com/Netcracker/qubership-core-lib-go-dbaas-postgres-client"
  "https://github.com/Netcracker/qubership-core-lib-go-error-handling"
  "https://github.com/Netcracker/qubership-core-lib-go-fiber-server-utils"
  "https://github.com/Netcracker/qubership-core-lib-go-maas-bg-segmentio"
  "https://github.com/Netcracker/qubership-core-lib-go-maas-client"
  "https://github.com/Netcracker/qubership-core-lib-go-maas-core"
  "https://github.com/Netcracker/qubership-core-lib-go-maas-segmentio"
  "https://github.com/Netcracker/qubership-core-lib-go-paas-mediation-client"
  "https://github.com/Netcracker/qubership-core-lib-go-rest-utils"
  "https://github.com/Netcracker/qubership-core-lib-go-stomp-websocket"
)

# ---------------------------------
# Requirements check
# ---------------------------------
need_bin() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "[ERROR] '$1' is required but not installed" >&2
    exit 1
  }
}
need_bin bash
need_bin git
need_bin sort
need_bin grep
need_bin sed
need_bin awk

# Absolute path for output file (next to the script)
OUTFILE="$(pwd)/output.yaml"
# Write YAML header once
cat > "$OUTFILE" <<'EOL'
go:
  libs-to-replace:
EOL

# ---------------------------------
# Working directory
# ---------------------------------
WORKDIR="$(mktemp -d)"
trap 'rm -rf "$WORKDIR"' EXIT
cd "$WORKDIR"

# Return latest semver tag from remote via git
# Accepts tags like: v1.2.3, 1.2.3, 1.2.3-RC1, 1.2.3+build.7
get_latest_semver_tag() {
  local repo_url="$1"

  git ls-remote --tags --quiet "$repo_url" \
  | awk '{
      ref=$2
      sub(/^refs\/tags\//,"",ref)   # drop "refs/tags/"
      sub(/\^\{\}$/,"",ref)         # drop peeled suffix "^{}"
      print ref
    }' \
  | sort -u \
  | grep -E '^v?[0-9]+\.[0-9]+\.[0-9]+([-.+][0-9A-Za-z.-]+)?$' \
  | sort -V \
  | tail -n 1
}

for repo_url in "${REPOS[@]}"; do
  repo_name="${repo_url##*/}"       # last path component

  echo "[INFO] Processing: $repo_url"
  echo "[INFO] Resolving latest semver tag via git ls-remote..."

  last_tag="$(get_latest_semver_tag "$repo_url" || true)"
  if [[ -z "${last_tag:-}" ]]; then
    echo "[ERROR] No semver tags found on remote: $repo_url" >&2
    exit 1
  fi

  echo "[INFO] Cloning at tag $last_tag ..."
  git -c advice.detachedHead=false clone --quiet --branch "$last_tag" --depth=1 "$repo_url" "$repo_name"
  cd "$repo_name"

  # Read module name from go.mod
  module_name=$(sed -n 's/^module[[:space:]]\+//p' go.mod)
  echo "[INFO] Module name: $module_name"

  {
    printf "    # %s@%s\n" "$repo_name" "$last_tag"
    cat <<EOL
    - old-name: $module_name
      new-name: $module_name
      new-version: $last_tag
EOL
  } >> "$OUTFILE"

  echo "[INFO] Done: $repo_name@$last_tag"
  echo ""
  cd ..
done

echo "[INFO] All done. Results saved to $OUTFILE"

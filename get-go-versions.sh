#!/bin/sh
set -eu

# List of test repositories
REPOS="
https://github.com/Netcracker/qubership-core-lib-go
https://github.com/Netcracker/qubership-core-lib-go-actuator-common
https://github.com/Netcracker/qubership-core-lib-go-bg-kafka
https://github.com/Netcracker/qubership-core-lib-go-bg-state-monitor
https://github.com/Netcracker/qubership-core-lib-go-dbaas-arangodb-client
https://github.com/Netcracker/qubership-core-lib-go-dbaas-base-client
https://github.com/Netcracker/qubership-core-lib-go-dbaas-cassandra-client
https://github.com/Netcracker/qubership-core-lib-go-dbaas-clickhouse-client
https://github.com/Netcracker/qubership-core-lib-go-dbaas-mongo-client
https://github.com/Netcracker/qubership-core-lib-go-dbaas-opensearch-client
https://github.com/Netcracker/qubership-core-lib-go-dbaas-postgres-client
https://github.com/Netcracker/qubership-core-lib-go-error-handling
https://github.com/Netcracker/qubership-core-lib-go-fiber-server-utils
https://github.com/Netcracker/qubership-core-lib-go-maas-bg-segmentio
https://github.com/Netcracker/qubership-core-lib-go-maas-client
https://github.com/Netcracker/qubership-core-lib-go-maas-core
https://github.com/Netcracker/qubership-core-lib-go-maas-segmentio
https://github.com/Netcracker/qubership-core-lib-go-paas-mediation-client
https://github.com/Netcracker/qubership-core-lib-go-rest-utils
https://github.com/Netcracker/qubership-core-lib-go-stomp-websocket

"

# Temporary working directory
WORKDIR="$(mktemp -d)"
trap 'rm -rf "$WORKDIR"' EXIT

# Absolute path for output file (next to the script)
OUTFILE="$(pwd)/out.txt"
: > "$OUTFILE"

cd "$WORKDIR"

{
	printf "go:\n"
	printf "  libs-to-replace:\n"
} >> "$OUTFILE"

for repo in $REPOS; do
    name=$(basename "$repo" .git)

    echo "[INFO] Initializing repo $repo ..."
    git init -q "$name"
    cd "$name"
    git remote add origin "$repo"

    echo "[INFO] Fetching tags for $name ..."
    git fetch --tags --quiet --depth=1 origin

    # Try to get the latest annotated tag
    if last_tag=$(git describe --tags --abbrev=0 2>/dev/null); then
        echo "[INFO] Found annotated tag: $last_tag"
    else
        echo "[WARN] No annotated tags found, falling back to lightweight tags..."
        last_tag=$(git tag --list | sort -V | tail -n 1) || true
        if [ -z "$last_tag" ]; then
            echo "[ERROR] No tags found at all in $repo" >&2
            exit 1
        fi
        echo "[INFO] Found lightweight tag: $last_tag"
    fi

    echo "[INFO] Checking out go.mod from $last_tag ..."
    git fetch --quiet --depth=1 origin "refs/tags/$last_tag"
    git checkout --quiet FETCH_HEAD -- go.mod

    # Read module name from go.mod
    module_name=$(sed -n 's/^module[[:space:]]\+//p' go.mod)
    echo "[INFO] Module name: $module_name"

    {
        printf "    - old-name: %s\n" "$module_name"
        printf "      new-name: %s\n" "$module_name"
        printf "      new-version: %s\n" "$last_tag"
    } >> "$OUTFILE"

    echo "[INFO] Wrote entry to $OUTFILE"
    cd ..
done

echo "[INFO] All done. Results saved to $OUTFILE"

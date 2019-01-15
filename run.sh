#!/usr/bin/env bash

build() {
    go clean -x

    go generate ./...

    go vet -all ./...

    [ -n "$SHADOW" ] && go vet -all "$SHADOW" ./... && echo "Shadowed-variables check complete."

    go build
}

if which shadow >/dev/null 2>/dev/null; then
    # Install via: go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
    SHADOW="-vettool=$(which shadow)"
fi

set -e  # Exit on error.

build

./joker -e '(print "\nLibraries available in this build:\n  ") (loaded-libs) (println)'

SUM256="$(go run tools/sum256dir/main.go std)"

(cd std; ../joker generate-std.joke)

NEW_SUM256="$(go run tools/sum256dir/main.go std)"

if [ "$SUM256" != "$NEW_SUM256" ]; then
    cat <<EOF
Rebuilding Joker, as the libraries have changed; then regenerating docs.

EOF
    build

    (cd docs; ../joker generate-docs.joke)
fi

./joker "$@"

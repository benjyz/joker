#!/usr/bin/env bash

[ -z "$OPTIMIZE_STARTUP" ] && OPTIMIZE_STARTUP=$([ -f NO-OPTIMIZE-STARTUP.flag ] && echo false || echo true)

build() {
  go clean
  rm -f core/a_*.go  # In case switching from a gen-code branch or similar (any existing files might break the build here)
  go generate ./...
  go vet -tags gen_code main.go repl.go
  go vet -tags gen_code ./core/... ./std/...
  (cd core; go fmt a_*.go > /dev/null)
  go build
  if $OPTIMIZE_STARTUP; then
      mv -f joker joker.slow
      go build -tags fast_init
      ln -f joker joker.fast
      echo "...built both joker.slow and joker.fast (aka joker)."
  else
      ln -f joker joker.slow
  fi
}

set -e  # Exit on error.

[ ! -f NO-GOSTD.flag ] && (cd tools/gostd && go build .) && ./tools/gostd/gostd --replace --joker .

build

if [ "$1" == "-v" ]; then
  ./joker -e '(print "\nLibraries available in this build:\n  ") (loaded-libs) (println)'
fi

# Check for changes in std, and run just-built Joker, only when building for host os/architecture.
SUM256="$(go run tools/sum256dir/main.go std)"
if [ ! -f NO-GEN.flag ]; then
    OUT="$(cd std; ../joker generate-std.joke 2>&1 | grep -v 'WARNING:.*already refers' | grep '.')" || : # grep returns non-zero if no lines match
    if [ -n "$OUT" ]; then
        echo "$OUT"
        echo >&2 "Unable to generate fresh library files; exiting."
        exit 2
    fi
fi
(cd std; go fmt ./... > /dev/null)
NEW_SUM256="$(go run tools/sum256dir/main.go std)"

if [ "$SUM256" != "$NEW_SUM256" ]; then
    echo 'std has changed, rebuilding...'
    build
    (cd docs; ../joker generate-docs.joke)
fi

./joker "$@"

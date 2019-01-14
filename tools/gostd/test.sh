#!/bin/bash

export GOENV="_tests/gold/$(go env GOARCH)-$(go env GOOS)"
mkdir -p "$GOENV"

EXIT="exit 99"
if [ "$1" = "--on-error" ]; then
    EXIT="$2"
fi

RC=0

rm -fr $GOENV/joker
mkdir -p $GOENV/joker/{core/data,std}
cp -pv ../../custom.go $GOENV/joker/
cp -pv ../../core/data/core.joke $GOENV/joker/core/data/
cp -pv ../../std/generate-custom.joke $GOENV/joker/std/

[ ! -x gostd ] && echo >&2 "No executable to test." && exit 99

./gostd --no-timestamp --output-code --verbose --go _tests/small 2>&1 | grep -v '^Default context:' > $GOENV/small.gold
git diff --quiet -u $GOENV/small.gold || { echo >&2 "FAILED: small test"; RC=1; $EXIT; }

./gostd --no-timestamp --output-code --verbose --go _tests/big --replace --joker $GOENV/joker 2>&1 | grep -v '^Default context:' > $GOENV/big.gold
git diff --quiet -u $GOENV/big.gold || { echo >&2 "FAILED: big test"; RC=1; $EXIT; }

./gostd --no-timestamp --output-code --verbose --go "$GOSRC" 2>&1 | grep -v '^Default context:' > $GOENV/gosrc.gold
git diff --quiet -u $GOENV/gosrc.gold || { echo >&2 "FAILED: \$GOSRC test"; RC=1; $EXIT; }

exit $RC

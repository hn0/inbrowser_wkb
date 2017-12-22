#!/usr/bin/bash

which emsdk >> /dev/null 2>&1

if [[ "$?" -eq 1 ]]; then
    source /opt/emsdk/emsdk_env.sh --build=Release
fi

emcc wkb.c -s WASM=1 -s EXPORTED_FUNCTIONS='["_convert"]' -s EXTRA_EXPORTED_RUNTIME_METHODS='["ccall"]' -o ../client/js/wkb_asm.js


# quick and dirty path update
sed -i -e 's+wkb_asm.wasm+js/wkb_asm.wasm+g' ../client/js/wkb_asm.js

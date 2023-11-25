#!/bin/bash

env GOOS=js GOARCH=wasm go build -o web/stealfromsanta.wasm fisherevans.com/stealingfromsanta

cp $(go env GOROOT)/misc/wasm/wasm_exec.js web/.

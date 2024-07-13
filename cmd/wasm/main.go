//go:build js && wasm

package main

import (
	appwasm "github.com/gjbae1212/hit-counter/internal/app/wasm"
)

func main() {
	appwasm.Run()
}

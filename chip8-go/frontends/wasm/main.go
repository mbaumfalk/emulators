package main

import (
	"bytes"
	"syscall/js"

	chip8 "github.com/mbaumfalk/emulators/chip8-go"
)

type canvas struct {
	global js.Value
}

var c canvas

func main() {
	c.global = js.Global()
}

func (c canvas) Clear() {
	c.global.Call("clear")
}

func (c canvas) Draw(x, y, data uint8) bool {
	return c.global.Call("draw", x, y, data).Bool()
}

var core *chip8.Chip8
var loaded = false
var i = 0

//export loadRom
func loadRom() {
	// TinyGo (as of 0.13.1) can not pass a js.Value or slice to an exported function
	// Go apparently has to call a global function instead
	value := js.Global().Get("rom")
	rom := make([]byte, value.Get("byteLength").Int())
	js.CopyBytesToGo(rom, value)

	var err error
	core, err = chip8.Init(bytes.NewReader(rom), c)
	if err != nil {
		panic(err)
	}
	loaded = true
}

//export runFrame
func runFrame() bool {
	if !loaded {
		return false
	}
	for ; i < 8; i++ {
		if !core.RunInstruction() {
			break
		}
	}
	i = 0
	core.Tick()
	return core.ST != 0
}

//export step
func step() {
	if !loaded {
		return
	}
	core.RunInstruction()
	i++
	if i == 8 {
		i = 0
		core.Tick()
	}
}

//export setKey
func setKey(key uint8, down bool) {
	if core == nil {
		return
	}
	value := uint16(1 << key)
	if down {
		core.Buttons |= value
	} else {
		core.Buttons &= ^value
	}
}

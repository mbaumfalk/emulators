// A frontend to the emulator which outputs to a terminal using ASCII escape codes
// This was used primarily for testing the emulator before making the WebAssembly frontend
// This frontend has no functioning for inputs or outputing sound

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	chip8 "github.com/mbaumfalk/emulators/chip8-go"
)

type terminal [32][8]uint8

func (t *terminal) Clear() {
	fmt.Print("\x1b[2J\x1b[1;1f")
	for y := 0; y < 32; y++ {
		t[y] = [8]uint8{}
		fmt.Println("\x1b[40m                                                                                                                                \x1b[0m")
	}
}

func (t *terminal) drawByte(x, y uint8) {
	fmt.Printf("\x1b[%d;%df", y+1, (x<<1)+1)
	offset := x % 8
	for n := uint8(0); n < 8; n++ {
		index := (x + n) >> 3
		if index > 7 {
			index = 0
			x -= 32
			y++
			if y == 32 {
				break
			}
			fmt.Println("\x1b[0m")
		}
		black := (t[y][index]>>((15-n-offset)%8))&1 == 0
		if black {
			fmt.Print("\x1b[40m")
		} else {
			fmt.Print("\x1b[107m")
		}
		fmt.Print("  ")
	}
	fmt.Print("\x1b[0m\x1b[33;1f")
}

func (t *terminal) Draw(x, y, data uint8) bool {
	if y >= 32 {
		return false
	}
	index := x >> 3
	offset := x % 8
	result := (t[y][index] & (data >> offset)) != 0
	t[y][index] ^= data >> offset

	y1 := y
	if index == 7 {
		y1++
		index = 0
	} else {
		index++
	}

	if y1 < 32 {
		result = result || (t[y1][index]&(data<<(8-offset))) != 0
		t[y1][index] ^= data << (8 - offset)
	}
	t.drawByte(x, y)
	return result
}

func main() {
	filename := flag.String("file", "", "Filename")
	disassemble := flag.Bool("dis", false, "Disassemble")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if *disassemble {
		rom := [0xe00]uint8{}
		n, err := file.Read(rom[:])
		if err != nil {
			panic(err)
		}
		for i := 0; i < n; i += 2 {
			fmt.Printf("0x%03X   %s\n", i+0x200, chip8.StringInstruction(rom[i], rom[i+1]))
		}
		return
	}

	t := &terminal{}
	core, err := chip8.Init(file, t)
	if err != nil {
		panic(err)
	}

	i := 0
	for {
		if i == 8 {
			i = 0
			core.Tick()
			time.Sleep(time.Second / 60)
		}
		core.RunInstruction()
		i++
	}
}

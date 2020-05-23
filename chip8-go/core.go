package chip8

import (
	"io"
	"math/rand"
	"time"
)

var digitSprites = [...]uint8{
	0xf0, 0x90, 0x90, 0x90, 0xf0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xf0, 0x10, 0xf0, 0x80, 0xf0, // 2
	0xf0, 0x10, 0xf0, 0x10, 0xf0, // 3
	0x90, 0x90, 0xf0, 0x10, 0x10, // 4
	0xf0, 0x80, 0xf0, 0x10, 0xf0, // 5
	0xf0, 0x80, 0xf0, 0x90, 0xf0, // 6
	0xf0, 0x10, 0x20, 0x40, 0x40, // 7
	0xf0, 0x90, 0xf0, 0x90, 0xf0, // 8
	0xf0, 0x90, 0xf0, 0x10, 0xf0, // 9
	0xf0, 0x90, 0xf0, 0x90, 0x90, // A
	0xe0, 0x90, 0xe0, 0x90, 0xe0, // B
	0xf0, 0x80, 0x80, 0x80, 0xf0, // C
	0xe0, 0x90, 0x90, 0x90, 0xe0, // D
	0xf0, 0x80, 0xf0, 0x80, 0xf0, // E
	0xf0, 0x80, 0xf0, 0x80, 0x80, // F
}

type Display interface {
	Clear()
	Draw(uint8, uint8, uint8) bool
}

type Chip8 struct {
	Memory []uint8
	Stack  [0x10]uint16

	V  [0x10]uint8
	PC uint16
	SP uint16
	I  uint16
	DT uint8
	ST uint8

	Buttons uint16

	rand    *rand.Rand
	display Display
}

func (c *Chip8) timerLoop() {
	for {
		if c.DT > 0 {
			c.DT--
		}
		if c.ST > 0 {
			c.ST--
		}
		time.Sleep(time.Second / 60)
	}
}

func (c *Chip8) Tick() {
	if c.DT > 0 {
		c.DT--
	}
	if c.ST > 0 {
		c.ST--
	}
}

func Init(reader io.Reader, display Display) (*Chip8, error) {
	c := Chip8{
		PC:      0x200,
		Memory:  make([]uint8, 0x1000),
		display: display,
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	copy(c.Memory, digitSprites[:])
	if _, err := reader.Read(c.Memory[0x200:]); err != nil {
		return nil, err
	}
	display.Clear()

	return &c, nil
}

func (c *Chip8) RunInstruction() bool {
	high := c.Memory[c.PC]
	low := c.Memory[c.PC+1]
	c.PC += 2

	highReg := high & 0xf
	lowReg := low >> 4
	switch high >> 4 {
	case 0x0:
		if high == 0x00 && low == 0xe0 {
			c.display.Clear()
		} else if high == 0x00 && low == 0xee {
			c.SP--
			c.PC = c.Stack[c.SP]
		}
	case 0x1:
		c.PC = uint16(highReg)<<8 | uint16(low)
	case 0x2:
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = uint16(highReg)<<8 | uint16(low)
	case 0x3:
		if c.V[highReg] == low {
			c.PC += 2
		}
	case 0x4:
		if c.V[highReg] != low {
			c.PC += 2
		}
	case 0x5:
		if c.V[highReg] == c.V[lowReg] {
			c.PC += 2
		}
	case 0x6:
		c.V[highReg] = low
	case 0x7:
		c.V[highReg] += low
	case 0x8:
		switch low & 0xf {
		case 0x0:
			c.V[highReg] = c.V[lowReg]
		case 0x1:
			c.V[highReg] |= c.V[lowReg]
		case 0x2:
			c.V[highReg] &= c.V[lowReg]
		case 0x3:
			c.V[highReg] ^= c.V[lowReg]
		case 0x4:
			c.V[highReg] += c.V[lowReg]
			if c.V[highReg] < c.V[lowReg] {
				c.V[0xf] = 1
			} else {
				c.V[0xf] = 0
			}
		case 0x5:
			if c.V[highReg] > c.V[lowReg] {
				c.V[0xf] = 1
			} else {
				c.V[0xf] = 0
			}
			c.V[highReg] -= c.V[lowReg]
		case 0x6:
			c.V[0xf] = c.V[highReg] & 1
			c.V[highReg] >>= 1
		case 0x7:
			if c.V[highReg] < c.V[lowReg] {
				c.V[0xf] = 1
			} else {
				c.V[0xf] = 0
			}
			c.V[highReg] = c.V[lowReg] - c.V[highReg]
		case 0xe:
			c.V[0xf] = c.V[highReg] >> 7
			c.V[highReg] <<= 1
		}
	case 0x9:
		if c.V[highReg] != c.V[lowReg] {
			c.PC += 2
		}
	case 0xa:
		c.I = uint16(highReg)<<8 | uint16(low)
	case 0xb:
		c.PC = (uint16(highReg)<<8 | uint16(low)) + uint16(c.V[0x0])
	case 0xc:
		c.V[highReg] = uint8(c.rand.Intn(0x100)) & low
	case 0xd:
		c.V[0xf] = 0
		for spriteIndex := uint(0); spriteIndex < uint(low&0xf); spriteIndex++ {
			if c.display.Draw(c.V[highReg], c.V[lowReg]+uint8(spriteIndex), c.Memory[uint(c.I)+spriteIndex]) {
				c.V[0xf] = 1
			}
		}
	case 0xe:
		switch low {
		case 0x9e:
			// TODO
			if (c.Buttons>>c.V[highReg])&1 == 1 {
				c.PC += 2
			}
		case 0xa1:
			if (c.Buttons>>c.V[highReg])&1 == 0 {
				c.PC += 2
			}
		}
	case 0xf:
		switch low {
		case 0x07:
			c.V[highReg] = c.DT
		case 0x0a:
			if c.Buttons == 0 {
				c.PC -= 2
				return false
			}
			for i := uint8(0); i < 16; i++ {
				if (c.Buttons>>i)&1 == 1 {
					c.V[highReg] = i
					break
				}
			}
		case 0x15:
			c.DT = c.V[highReg]
		case 0x18:
			c.ST = c.V[highReg]
		case 0x1e:
			c.I += uint16(c.V[highReg])
		case 0x29:
			c.I = 5 * uint16(c.V[highReg])
		case 0x33:
			value := c.V[highReg]
			c.Memory[c.I] = value / 100
			c.Memory[c.I+1] = (value / 10) % 10
			c.Memory[c.I+2] = value % 10
		case 0x55:
			for i := uint16(0); i <= uint16(highReg); i++ {
				c.Memory[c.I+i] = c.V[i]
			}
		case 0x65:
			for i := uint16(0); i <= uint16(highReg); i++ {
				c.V[i] = c.Memory[c.I+i]
			}
		}
	}
	return true
}

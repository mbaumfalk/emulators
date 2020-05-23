package chip8

import "fmt"

func StringInstruction(high, low byte) string {
	highReg := high & 0xf
	switch high >> 4 {
	case 0x0:
		if high == 0x00 && low == 0xe0 {
			return "CLS"
		} else if high == 0x00 && low == 0xee {
			return "RET"
		} else {
			return fmt.Sprintf("SYS %X%02X (!)", highReg, low)
		}
	case 0x1:
		return fmt.Sprintf("JP $%X%02X", highReg, low)
	case 0x2:
		return fmt.Sprintf("CALL $%X%02X", highReg, low)
	case 0x3:
		return fmt.Sprintf("SE V%X, $%02X", highReg, low)
	case 0x4:
		return fmt.Sprintf("SNE V%X, $%02X", highReg, low)
	case 0x5:
		return fmt.Sprintf("SE V%X, V%X", highReg, low&0xf)
	case 0x6:
		return fmt.Sprintf("LD V%X, $%02X", highReg, low)
	case 0x7:
		return fmt.Sprintf("ADD V%X, $%02X", highReg, low)
	case 0x8:
		switch low & 0xf {
		case 0x0:
			return fmt.Sprintf("LD V%X, V%X", highReg, low>>4)
		case 0x1:
			return fmt.Sprintf("OR V%X, V%X", highReg, low>>4)
		case 0x2:
			return fmt.Sprintf("AND V%X, V%X", highReg, low>>4)
		case 0x3:
			return fmt.Sprintf("XOR V%X, V%X", highReg, low>>4)
		case 0x4:
			return fmt.Sprintf("ADD V%X, V%X", highReg, low>>4)
		case 0x5:
			return fmt.Sprintf("SUB V%X, V%X", highReg, low>>4)
		case 0x6:
			return fmt.Sprintf("SHR V%X, V%X", highReg, low>>4)
		case 0x7:
			return fmt.Sprintf("SUBN V%X, V%X", highReg, low>>4)
		case 0xe:
			return fmt.Sprintf("SHL V%X, V%X", highReg, low>>4)
		}
	case 0x9:
		return fmt.Sprintf("SNE V%X, V%X", highReg, low>>4)
	case 0xa:
		return fmt.Sprintf("LD I, $%X%02X", highReg, low)
	case 0xb:
		return fmt.Sprintf("JP V0, $%X%02X", highReg, low)
	case 0xc:
		return fmt.Sprintf("RND V%X, $%02X", highReg, low)
	case 0xd:
		return fmt.Sprintf("DRW V%X, V%X, %X", highReg, low>>4, low&0xf)
	case 0xe:
		switch low {
		case 0x9e:
			return fmt.Sprintf("SKP V%X", highReg)
		case 0xa1:
			return fmt.Sprintf("SKNP V%X", highReg)
		}
	case 0xf:
		switch low {
		case 0x07:
			return fmt.Sprintf("LD V%X, DT", highReg)
		case 0x0a:
			return fmt.Sprintf("LD V%X, K", highReg)
		case 0x15:
			return fmt.Sprintf("LD DT, V%X", highReg)
		case 0x18:
			return fmt.Sprintf("LD ST, V%X", highReg)
		case 0x1e:
			return fmt.Sprintf("ADD I< V%X", highReg)
		case 0x29:
			return fmt.Sprintf("LD F, V%X", highReg)
		case 0x33:
			return fmt.Sprintf("LD B, V%X", highReg)
		case 0x55:
			return fmt.Sprintf("LD [I], V%X", highReg)
		case 0x65:
			return fmt.Sprintf("LD V%X, [I]", highReg)
		}
	}
	// return "Unknown opcode"
	return fmt.Sprintf("$%02X%02X (!)", high, low)
}

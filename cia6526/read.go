package cia6526

import (
	"log"
	"os"
)

var count = 0
var loggerRead = false

func (C *CIA) dispReadReg(label string, mode byte, pr byte, key byte, res byte) {
	if C.name == "CIA1" && loggerRead {
		log.Printf("%-12s -*-    Mode: %08b -  RegVal: %08b -    Key: %08b - Res: %08b", label, C.Reg[mode], C.Reg[pr], key, res)
	}
}

func (C *CIA) stopOncount() {
	count++
	if count > 40 {
		os.Exit(1)
	}
}

func (C *CIA) Read(addr uint16) byte {
	reg := byte(addr) // - ((addr >> 4) << 4)
	// clog.Trace("CIA", "Read", "%s - addr: %04X - Reg: %02X (%d)", C.name, addr, reg, reg)
	switch reg {
	case PRA:
		if C.name == "CIA1" {
			res := byte(0xFF)
			if C.InputLine.KeyCode != 0 {
				C.buffer = keyMap[C.InputLine.Mode]
				if C.InputLine.Mode != 0 {
					if C.Reg[PRB]|C.buffer.row == C.buffer.row {
						res &= (^C.buffer.col & ^C.Reg[DDRA]) ^ C.Reg[PRA]
					}
				}
				C.buffer = keyMap[C.InputLine.KeyCode]
				if C.Reg[PRB]|C.buffer.row == C.buffer.row {
					res &= (^C.buffer.col & ^C.Reg[DDRA]) ^ C.Reg[PRA]
					C.dispReadReg("READ PRA", DDRA, PRA, C.buffer.col, res)
				}
				return res
			}
			return C.Reg[PRA] & ^C.Reg[DDRA]
		}
		// C.dispReadReg("READ PRA", DDRA, PRA, C.buffer.col, C.Reg[PRA] & ^C.Reg[DDRA])
		return C.Reg[PRA]

	case PRB:
		if C.name == "CIA1" && C.InputLine.KeyCode != 0 {
			res := byte(0xFF)
			if C.InputLine.KeyCode != 0 {
				C.buffer = keyMap[C.InputLine.Mode]
				if C.InputLine.Mode != 0 {
					if C.Reg[PRA]|C.buffer.col == C.buffer.col {
						res &= (^C.buffer.row & ^C.Reg[DDRB]) ^ C.Reg[PRB]
					}
				}
				C.buffer = keyMap[C.InputLine.KeyCode]
				if C.Reg[PRA]|C.buffer.col == C.buffer.col {
					res &= (^C.buffer.row & ^C.Reg[DDRB]) ^ C.Reg[PRB]
					C.dispReadReg("READ PRB", DDRB, PRB, C.buffer.row, res)
				}
				return res
			}
			return C.Reg[PRB] & ^C.Reg[DDRB]
		}
		// C.dispReadReg("READ PRB", DDRB, PRB, C.buffer.row, C.Reg[PRB] & ^C.Reg[DDRB])
		return C.Reg[PRB]

	case ICR:
		val := C.Reg[ICR]
		C.Reg[ICR] = 0
		*C.Signal_Pin = 0
		return val
	default:
		return C.Reg[reg]
	}
}

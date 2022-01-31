package cia6526

import (
	"log"
	"os"
)

var count = 0
var loggerRead = false

func (C *CIA) dispReadReg(label string, mode uint16, pr uint16, key byte, res byte) {
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
	reg := addr - ((addr >> 4) << 4)
	switch reg {
	case PRA:
		if C.name == "CIA1" && C.InputLine.KeyCode != 0 {
			res := byte(0xFF)
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
		// C.dispReadReg("READ PRA", DDRA, PRA, C.buffer.col, C.Reg[PRA] & ^C.Reg[DDRA])
		return C.Reg[PRA] & ^C.Reg[DDRA]

	case PRB:
		if C.name == "CIA1" && C.InputLine.KeyCode != 0 {
			res := byte(0xFF)
			// fmt.Printf("PRA: %08b - ", C.Reg[PRA])
			C.buffer = keyMap[C.InputLine.Mode]
			if C.InputLine.Mode != 0 {
				if C.Reg[PRA]|C.buffer.col == C.buffer.col {
					// fmt.Printf("Shift Col: %08b  - Shift Row: %08b - ", C.buffer.col, C.buffer.row)
					res &= (^C.buffer.row & ^C.Reg[DDRB]) ^ C.Reg[PRB]
					// fmt.Printf("Res1: %08b -*- ", res)
				}
			}
			C.buffer = keyMap[C.InputLine.KeyCode]
			if C.Reg[PRA]|C.buffer.col == C.buffer.col {
				// fmt.Printf("Key Col: %08b  - Key Row: %08b - ", C.buffer.col, C.buffer.row)
				res &= (^C.buffer.row & ^C.Reg[DDRB]) ^ C.Reg[PRB]
				// fmt.Printf("Res2: %08b - ", res)
				C.dispReadReg("READ PRB", DDRB, PRB, C.buffer.row, res)
			}
			// fmt.Printf("Final Res: %08b\n", res)
			return res
		}
		// C.dispReadReg("READ PRB", DDRB, PRB, C.buffer.row, C.Reg[PRB] & ^C.Reg[DDRB])
		return C.Reg[PRB] & ^C.Reg[DDRB]

	case ICR:
		val := C.Reg[ICR]
		C.Reg[ICR] = 0
		*C.Signal_Pin = 0
		return val
	default:
		return C.Reg[reg]
	}
}

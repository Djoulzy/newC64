package cia6526

import (
	"log"
	"os"
)

var count = 0
var logger = false

func (C *CIA) dispReadReg(label string, mode uint16, pr uint16, key byte, res byte) {
	if C.name == "CIA1" && logger {
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
		if C.name == "CIA1" && C.InputLine != 0 {
			C.buffer = keyMap[C.InputLine]
			if C.Reg[PRB]|C.buffer.row == C.buffer.row {
				test := ^C.buffer.col & ^C.Reg[DDRA]
				res := test ^ C.Reg[PRA]

				C.dispReadReg("READ PRA", DDRA, PRA, C.buffer.col, res)
				// C.stopOncount()
				return res
			}
		}
		C.dispReadReg("READ PRA", DDRA, PRA, C.buffer.col, C.Reg[PRA] & ^C.Reg[DDRA])
		return C.Reg[PRA] & ^C.Reg[DDRA]

	case PRB:
		if C.name == "CIA1" && C.InputLine != 0 {
			C.buffer = keyMap[C.InputLine]
			if C.Reg[PRA]|C.buffer.col == C.buffer.col {
				test := ^C.buffer.row & ^C.Reg[DDRB]
				res := test ^ C.Reg[PRB]

				C.dispReadReg("READ PRB", DDRB, PRB, C.buffer.row, res)
				// C.stopOncount()
				return res
			}
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

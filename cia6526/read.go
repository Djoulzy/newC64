package cia6526

import (
	"log"
	"os"
)

var count = 0

var buffer = Keyb_B

func (C *CIA) dispReadReg(label string, mode uint16, pr uint16, key byte, res byte) {
	if C.name == "CIA1" {
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
		keyLevel := byte(0xFF)
		if C.Reg[PRB]&buffer.row > 0 {
			keyLevel = 0x00
		}
		val := buffer.col & keyLevel
		test := val & ^C.Reg[DDRA]
		res := test ^ C.Reg[PRA]
		// C.dispReadReg("READ PRA", DDRA, PRA, buffer.col, res)
		// C.stopOncount()
		return res
	case PRB:
		keyLevel := byte(0xFF)
		if C.Reg[PRA]&buffer.col > 0 {
			keyLevel = 0x00
		}
		val := buffer.row & keyLevel
		test := val & ^C.Reg[DDRB]
		res := test ^ C.Reg[PRB]
		// C.dispReadReg("READ PRB", DDRB, PRB, buffer.row, res)
		// C.stopOncount()
		return res
	case ICR:
		val := C.Reg[ICR]
		C.Reg[ICR] = 0
		*C.Signal_Pin = 0
		return val
	default:
		return C.Reg[reg]
	}
}

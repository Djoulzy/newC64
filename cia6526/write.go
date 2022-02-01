package cia6526

import "log"

var loggerWrite bool = false

func (C *CIA) dispWriteReg(label string, mode uint16, pr uint16, val byte, res byte) {
	if C.name == "CIA1" && loggerWrite {
		log.Printf("%-12s -*-    Mode: %08b -  RegVal: %08b - NewVal: %08b - Res: %08b", label, C.Reg[mode], C.Reg[pr], val, res)
	}
}

func (C *CIA) dispWriteDir(label string, mode uint16, pr uint16, val byte, res byte) {
	if C.name == "CIA1" && loggerWrite {
		log.Printf("%-12s -*- NewMode: %08b -> RegVal: %08b - NewVal: %08b", label, C.Reg[mode], C.Reg[pr], res)
	}
}

func (C *CIA) Write(addr uint16, val byte) {

	reg := addr - ((addr >> 4) << 4)

	switch reg {
	case PRA:
		newPr := (C.Reg[PRA] & ^C.Reg[DDRA]) | (val & C.Reg[DDRA])
		if C.buffer != Keyb_NULL {
			C.dispWriteReg("WRITE PRA", DDRA, PRA, val, newPr)
		}
		C.Reg[PRA] = newPr
		if C.name == "CIA2" {
			log.Printf("Mode %08b - Val %08b", C.Reg[DDRA], val)
			*C.VICBankSelect = int(C.Reg[PRA] & 0b00000011)
		}
	case PRB:
		newPr := (C.Reg[PRB] & ^C.Reg[DDRB]) | (val & C.Reg[DDRB])
		// if C.buffer != Keyb_NULL {
		// 	C.dispWriteReg("WRITE PRB", DDRB, PRB, val, newPr)
		// }
		C.Reg[PRB] = newPr
	case DDRA:
		C.Reg[DDRA] = val
		newPr := C.Reg[PRA] | ^val
		// if C.buffer != Keyb_NULL {
		// 	C.dispWriteDir("WRITE DDRA", DDRA, PRA, val, newPr)
		// }
		C.Reg[PRA] = newPr
	case DDRB:
		C.Reg[DDRB] = val
		newPr := C.Reg[PRB] | ^val
		// if C.buffer != Keyb_NULL {
		// 	C.dispWriteDir("WRITE DDRB", DDRB, PRB, val, newPr)
		// }
		C.Reg[PRB] = newPr
	case TALO:
		C.timerA_latchLO = val
	case TAHI:
		C.timerA_latchHI = val
		if C.Reg[CRA]&CTRL_START_STOP == 0 {
			C.Reg[TALO] = C.timerA_latchLO
			C.Reg[TAHI] = C.timerA_latchHI
		}
	case TBLO:
		C.timerB_latchLO = val
	case TBHI:
		C.timerB_latchHI = val
		if C.Reg[CRB]&CTRL_START_STOP == 0 {
			C.Reg[TBLO] = C.timerB_latchLO
			C.Reg[TBHI] = C.timerB_latchHI
		}
	case TOD10THS:
	case TODSEC:
	case TODMIN:
	case TODHR:
	case SRD:
	case ICR:
		mask := val & 0b00001111
		if mask > 0 {
			if val&0b10000000 > 0 { // 7eme bit = 1 -> mask set
				C.interrupt_mask = C.interrupt_mask | mask
			} else {
				C.interrupt_mask = C.interrupt_mask & ^mask
			}
		}
	case CRA:
		// Load Latch Once
		if val&CTRL_LOAD_LATCH > 0 {
			C.Reg[TAHI] = C.timerA_latchHI
			C.Reg[TALO] = C.timerA_latchLO
			C.cycleDelay = 2
		}
		C.Reg[CRA] = val & 0b11100111
	case CRB:
		// Load Latch Once
		if val&CTRL_LOAD_LATCH > 0 {
			C.Reg[TBHI] = C.timerB_latchHI
			C.Reg[TBLO] = C.timerB_latchLO
			C.cycleDelay = 2
		}
		C.Reg[CRB] = val & 0b11100111
	}
}

func (C *CIA) testBit(reg uint16, mask byte) bool {
	if C.Reg[reg]&mask == mask {
		return true
	}
	return false
}

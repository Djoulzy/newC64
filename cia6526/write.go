package cia6526

func (C *CIA) Write(addr uint16, val byte) {

	reg := addr - ((addr >> 4) << 4)

	switch reg {
	case PRA:
		C.Reg[PRA] = 0b10000001
	case PRB:
		C.Reg[PRB] = 0b11111111
	case DDRA:
		fallthrough
	case DDRB:
		C.Reg[reg] = val
	case TALO:
		C.timerA_latchLO = val
	case TAHI:
		C.timerA_latchHI = val
	case TBLO:
		C.timerB_latchLO = val
	case TBHI:
		C.timerB_latchHI = val
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
		}
		C.Reg[CRA] = val & 0b11100111
	case CRB:
		// Load Latch Once
		if val&CTRL_LOAD_LATCH > 0 {
			C.Reg[TBHI] = C.timerB_latchHI
			C.Reg[TBLO] = C.timerB_latchLO
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

package cia6526

func (C *CIA) Write(addr uint16, val byte) {

	reg := addr - ((addr >> 4) << 4)

	switch reg {
	case PRA:
	case PRB:
	case DDRA:
	case DDRB:
	case TALO:
		C.timerAlatch = int32(C.Reg[TAHI])<<8 + int32(val)
	case TAHI:
		C.timerAlatch = int32(val)<<8 + int32(C.Reg[TALO])
	case TBLO:
		C.timerBlatch = int32(C.Reg[TBHI])<<8 + int32(val)
	case TBHI:
		C.timerBlatch = int32(val)<<8 + int32(C.Reg[TBLO])
	case TOD10THS:
	case TODSEC:
	case TODMIN:
	case TODHR:
	case SRD:
	case ICR:
		mask := val & 0b00001111
		if mask > 0 {
			if val&0b10000000 > 0 { // 7eme bit = 1 -> mask set
				C.Reg[ICR] = C.Reg[ICR] | mask
			} else {
				C.Reg[ICR] = C.Reg[ICR] & ^mask
			}
		}
	case CRA:
		// Load Latch Once
		if val&0b00010000 > 0 {
			C.timerAlatch = int32(C.Reg[TAHI])<<8 + int32(C.Reg[TALO])
		}
		// Start or stop timer
		if val&0b00000001 == 1 {
			C.timerAstate = true
		} else {
			C.timerAstate = false
		}
		C.Reg[CRA] = val & 0b11101111
	case CRB:
		// Load Latch Once
		if val&0b00010000 > 0 {
			C.timerAlatch = int32(C.Reg[TBHI])<<8 + int32(C.Reg[TBLO])
		}
		// Start or stop timer
		if val&0b00000001 == 1 {
			C.timerBstate = true
		} else {
			C.timerBstate = false
		}
		C.Reg[CRB] = val & 0b11101111
	}
}

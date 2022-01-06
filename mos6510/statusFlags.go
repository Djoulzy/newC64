package mos6510

type toggle byte

func (C *CPU) updateZ(val byte) {
	if val == 0 {
		C.S |= ^Z_mask
	} else {
		C.S &= Z_mask
	}
}

func (C *CPU) updateN(val byte) {
	if val&0b10000000 > 0 {
		C.S |= ^N_mask
	} else {
		C.S &= N_mask
	}
}

func (C *CPU) setN(val bool) {
	if val {
		C.S |= ^N_mask
	} else {
		C.S &= N_mask
	}
}

func (C *CPU) setC(val bool) {
	if val {
		C.S |= ^C_mask
	} else {
		C.S &= C_mask
	}
}

func (C *CPU) issetC() bool {
	return C.S & ^C_mask > 0
}

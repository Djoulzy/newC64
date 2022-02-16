package mos6510

import "log"

//////////////////////////////////
//////////// Zero Flag ///////////
//////////////////////////////////

func (C *CPU) updateZ(val byte) {
	if val == 0 {
		C.S |= ^Z_mask
	} else {
		C.S &= Z_mask
	}
}

func (C *CPU) setZ(val bool) {
	if val {
		C.S |= ^Z_mask
	} else {
		C.S &= Z_mask
	}
}

func (C *CPU) issetZ() bool {
	return C.S & ^Z_mask > 0
}

//////////////////////////////////
///////// Negative Flag //////////
//////////////////////////////////

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

func (C *CPU) issetN() bool {
	return C.S & ^N_mask > 0
}

//////////////////////////////////
/////////// Carry Flag ///////////
//////////////////////////////////

func (C *CPU) setC(val bool) {
	if val {
		C.S |= ^C_mask
	} else {
		C.S &= C_mask
	}
}

func (C *CPU) getC() byte {
	if C.S & ^C_mask > 0 {
		return 0x01
	}
	return 0x00
}

func (C *CPU) issetC() bool {
	return C.S & ^C_mask > 0
}

//////////////////////////////////
/////////// Break Flag ///////////
//////////////////////////////////

func (C *CPU) setB(val bool) {
	if val {
		C.S |= ^B_mask
	} else {
		C.S &= B_mask
	}
	// log.Printf("Break flag: %04X - %s", C.InstStart, C.registers())
}

//////////////////////////////////
////////// Unused Flag ///////////
//////////////////////////////////

func (C *CPU) setU(val bool) {
	if val {
		C.S |= ^U_mask
	} else {
		C.S &= U_mask
	}
}

//////////////////////////////////
///////// Overflow Flag //////////
//////////////////////////////////

func (C *CPU) updateV(m, n, result byte) {
	if (m^result)&(n^result)&0x80 != 0 {
		C.S |= ^V_mask
	} else {
		C.S &= V_mask
	}
}

func (C *CPU) issetV() bool {
	return C.S & ^V_mask > 0
}

func (C *CPU) setV(val bool) {
	if val {
		C.S |= ^V_mask
	} else {
		C.S &= V_mask
	}
}

//////////////////////////////////
///////// Decimal Flag ///////////
//////////////////////////////////

func (C *CPU) setD(val bool) {
	if val {
		log.Printf("%04X - %s", C.InstStart, C.registers())
		// log.Fatal("Decimal flag")
		C.S |= ^D_mask
	} else {
		C.S &= D_mask
	}
}

func (C *CPU) getD() byte {
	if C.S & ^D_mask > 0 {
		return 0x01
	}
	return 0x00
}

//////////////////////////////////
//////// Interrupt Flag //////////
//////////////////////////////////

func (C *CPU) setI(val bool) {
	if val {
		C.S |= ^I_mask
	} else {
		C.S &= I_mask
	}
}

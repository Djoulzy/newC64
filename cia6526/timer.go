package cia6526

func (C *CIA) TimerA() {
	C.timerAlatch--
	C.Reg[TAHI]=byte(C.timerAlatch >> 8)
	C.Reg[TALO]=byte(C.timerAlatch)
	if C.timerAlatch < 0 {
		icr := C.Reg[ICR]
		// log.Println("underflow timer A")
		if (icr&0b00000001 > 0) && (icr&0b1000000 == 0) {
			C.Reg[ICR]=icr | 0b10000001
			// log.Printf("%s: Int timer A\n", C.name)
			*C.Signal_Pin = 1
		}
		if C.Reg[CRA]&0b00001000 > 0 {
			return
		} else {
			C.timerBlatch = int32(C.Reg[TAHI])<<8 + int32(C.Reg[TALO])
		}
	}
}

func (C *CIA) TimerB() {
	C.timerAlatch--
	C.Reg[TBHI]=byte(C.timerBlatch >> 8)
	C.Reg[TBLO]=byte(C.timerBlatch)
	if C.timerAlatch < 0 {
		icr := C.Reg[ICR]
		// log.Println("underflow timer B")
		if (icr&0b00000010 > 0) && (icr&0b1000000 == 0) {
			C.Reg[ICR]=icr | 0b10000010
			// log.Printf("%s: Int timer B\n", C.name)
			*C.Signal_Pin = 1
		}
		if C.Reg[CRB]&0b00001000 > 0 {
			return
		} else {
			C.timerBlatch = int32(C.Reg[TBHI])<<8 + int32(C.Reg[TBLO])
		}
	}
}

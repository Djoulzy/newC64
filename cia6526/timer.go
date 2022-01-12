package cia6526

func (C *CIA) TimerA() {
	C.timerAlatch--
	C.reg[TAHI].Output(byte(C.timerAlatch >> 8))
	C.reg[TALO].Output(byte(C.timerAlatch))
	if C.timerAlatch < 0 {
		icr := C.reg[ICR].Val
		// log.Println("underflow timer A")
		if (icr&0b00000001 > 0) && (icr&0b1000000 == 0) {
			C.reg[ICR].Output(icr | 0b10000001)
			// log.Printf("%s: Int timer A\n", C.name)
			*C.Signal_Pin = 1
		}
		if C.reg[CRA].Val&0b00001000 > 0 {
			return
		} else {
			C.timerBlatch = int32(C.reg[TAHI].Val)<<8 + int32(C.reg[TALO].Val)
		}
	}
}

func (C *CIA) TimerB() {
	C.timerAlatch--
	C.reg[TBHI].Output(byte(C.timerBlatch >> 8))
	C.reg[TBLO].Output(byte(C.timerBlatch))
	if C.timerAlatch < 0 {
		icr := C.reg[ICR].Val
		// log.Println("underflow timer B")
		if (icr&0b00000010 > 0) && (icr&0b1000000 == 0) {
			C.reg[ICR].Output(icr | 0b10000010)
			// log.Printf("%s: Int timer B\n", C.name)
			*C.Signal_Pin = 1
		}
		if C.reg[CRB].Val&0b00001000 > 0 {
			return
		} else {
			C.timerBlatch = int32(C.reg[TBHI].Val)<<8 + int32(C.reg[TBLO].Val)
		}
	}
}

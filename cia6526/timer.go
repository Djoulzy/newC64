package cia6526

import "newC64/memory"

func (C *CIA) TimerA() {
	C.timerAlatch--
	C.io.CiaRegWrite(TAHI, byte(C.timerAlatch>>8), memory.NONE)
	C.io.CiaRegWrite(TALO, byte(C.timerAlatch), memory.NONE)
	if C.timerAlatch < 0 {
		// log.Println("underflow timer A")
		if (C.io.Val[ICR]&0b00000001 > 0) && (C.io.Val[ICR]&0b1000000 == 0) {
			C.io.CiaRegWrite(ICR, C.io.Val[ICR]|0b10000001, memory.NONE)
			// log.Printf("%s: Int timer A\n", C.name)
			*C.Signal_Pin = 1
		}
		if C.io.Val[CRA]&0b00001000 > 0 {
			return
		} else {
			C.timerAlatch = int32(C.io.Val[TAHI])<<8 + int32(C.io.Val[TALO])
		}
	}
}

func (C *CIA) TimerB() {
	C.timerAlatch--
	C.io.CiaRegWrite(TBHI, byte(C.timerBlatch>>8), memory.NONE)
	C.io.CiaRegWrite(TBLO, byte(C.timerBlatch), memory.NONE)
	if C.timerAlatch < 0 {
		// log.Println("underflow timer B")
		if (C.io.Val[ICR]&0b00000010 > 0) && (C.io.Val[ICR]&0b1000000 == 0) {
			C.io.CiaRegWrite(ICR, C.io.Val[ICR]|0b10000010, memory.NONE)
			// log.Printf("%s: Int timer B\n", C.name)
			*C.Signal_Pin = 1
		}
		if C.io.Val[CRB]&0b00001000 > 0 {
			return
		} else {
			C.timerBlatch = int32(C.io.Val[TBHI])<<8 + int32(C.io.Val[TBLO])
		}
	}
}

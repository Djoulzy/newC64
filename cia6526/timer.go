package cia6526

func (C *CIA) TimerA() {
	last := *C.systemCycle
	for {
		select {
		case <-C.timerAcom:
			return
		default:
			if *C.systemCycle != last {
				C.timerAlatch--
				C.io.CiaRegWrite(TAHI, byte(C.timerAlatch>>8))
				C.io.CiaRegWrite(TALO, byte(C.timerAlatch))
				if C.timerAlatch < 0 {
					// log.Println("underflow timer A")
					if (C.io.Val[ICR]&0b00000001 > 0) && (C.io.Val[ICR]&0b1000000 == 0) {
						C.io.CiaRegWrite(ICR, C.io.Val[ICR]|0b10000001)
						// log.Printf("%s: Int timer A\n", C.name)
						*C.Signal_Pin = 1
					}
					if C.io.Val[CRA]&0b00001000 > 0 {
						return
					} else {
						C.timerAlatch = int32(C.io.Val[TAHI])<<8 + int32(C.io.Val[TALO])
					}
				}
				// log.Printf("timerA: %d\n", C.timerAlatch)
				last = *C.systemCycle
			}
		}
	}
}

func (C *CIA) TimerB() {
	last := *C.systemCycle
	for {
		select {
		case <-C.timerBcom:
			return
		default:
			if *C.systemCycle != last {
				C.timerAlatch--
				C.io.CiaRegWrite(TBHI, byte(C.timerBlatch>>8))
				C.io.CiaRegWrite(TBLO, byte(C.timerBlatch))
				if C.timerAlatch < 0 {
					// log.Println("underflow timer B")
					if (C.io.Val[ICR]&0b00000010 > 0) && (C.io.Val[ICR]&0b1000000 == 0) {
						C.io.CiaRegWrite(ICR, C.io.Val[ICR]|0b10000010)
						// log.Printf("%s: Int timer B\n", C.name)
						*C.Signal_Pin = 1
					}
					if C.io.Val[CRB]&0b00001000 > 0 {
						return
					} else {
						C.timerBlatch = int32(C.io.Val[TBHI])<<8 + int32(C.io.Val[TBLO])
					}
				}
				// log.Printf("timerA: %d\n", C.timerAlatch)
				last = *C.systemCycle
			}
		}
	}
}

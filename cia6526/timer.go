package cia6526

const (
	CTRL_START_STOP   byte = 0b00000001
	CTRL_LOOP         byte = 0b00000100
	CTRL_LOAD_LATCH   byte = 0b00010000
	CTRL_SET_CLK_FREQ byte = 0b10000000
)

func (C *CIA) time_is_over(lo *byte, hi *byte) bool {
	*lo--
	if *lo == 0 {
		if *hi == 0 {
			return true
		}
		*hi--
		*lo = 0xFF
	}
	return false
}

func (C *CIA) TimerA() {
	// log.Printf("Tick Timer A of %s", C.name)
	if C.time_is_over(&C.Reg[TALO], &C.Reg[TAHI]) {
		// log.Println("underflow timer A")
		if C.interrupt_mask&INT_UNDERFL_TA > 0 {
			C.Reg[ICR] |= INT_UNDERFL_TA | INT_SET
			// log.Printf("%s: Int timer A\n", C.name)
			*C.Signal_Pin = 1
		}
		if C.Reg[CRA]&CTRL_LOOP > 0 {
			C.Reg[CRA] &= ^CTRL_START_STOP
			return
		} else {
			C.Reg[TAHI] = C.timerA_latchHI
			C.Reg[TALO] = C.timerA_latchLO
		}
	}
}

func (C *CIA) TimerB() {
	// log.Printf("Tick Timer B of %s", C.name)
	if C.time_is_over(&C.Reg[TBLO], &C.Reg[TBHI]) {
		// log.Println("underflow timer B")
		if C.interrupt_mask&INT_UNDERFL_TB > 0 {
			C.Reg[ICR] |= INT_UNDERFL_TB | INT_SET
			// log.Printf("%s: Int timer B\n", C.name)
			*C.Signal_Pin = 1
		}
		if C.Reg[CRB]&CTRL_LOOP > 0 {
			C.Reg[CRB] &= ^CTRL_START_STOP
			return
		} else {
			C.Reg[TBHI] = C.timerB_latchHI
			C.Reg[TBLO] = C.timerB_latchLO
		}
	}
}

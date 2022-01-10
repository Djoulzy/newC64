package cia6526

import "newC64/memory"

type CIA struct {
	name        string
	io          *memory.MEM
	Signal_Pin  *int
	systemCycle *uint16

	timerAlatch int32
	timerAcom   chan int

	timerBlatch int32
	timerBcom   chan int
}

const (
	PRA uint16 = iota
	PRB
	DDRA
	DDRB
	TALO
	TAHI
	TBLO
	TBHI
	TOD10THS
	TODSEC
	TODMIN
	TODHR
	SRD
	ICR // Interrupt control register
	CRA // Timer A Control
	CRB // Timer B Control
)

func (C *CIA) Init(name string, memCells *memory.MEM, timer *uint16) {
	C.name = name
	C.io = memCells
	C.systemCycle = timer

	C.timerAlatch = 0
	C.timerAcom = make(chan int)
	C.timerBlatch = 0
	C.timerBcom = make(chan int)

	// C.SetValue(PRA, 0x0)
	// C.SetValue(PRB, 0xFF)
	// C.SetValue(DDRA, 0)
	// C.SetValue(DDRB, 0)
	// C.SetValue(TALO, 0xFF)
	// C.SetValue(TAHI, 0xFF)
	// C.SetValue(TBLO, 0xFF)
	// C.SetValue(TBHI, 0xFF)
	// C.SetValue(TOD10THS, 0)
	// C.SetValue(TODSEC, 0)
	// C.SetValue(TODMIN, 0)
	// C.SetValue(TODHR, 0x01)
	// C.SetValue(SRD, 0)
	// C.SetValue(ICR, 0)
	// C.mem[ICR].Zone[mem.RAM] = 0b00001111
	// C.SetValue(CRA, 0)
	// C.SetValue(CRB, 0)
}

// func (C *CIA) SetValue(port byte, value byte) {
// 	for i := 0; i < 16; i++ {
// 		zone := port + byte(16*i)
// 		C.mem[zone].Zone[mem.IO] = value
// 		C.mem[zone].Zone[mem.RAM] = value
// 	}
// }

// func (C *CIA) execTimerA() {
// 	reg := C.mem[CRA].Zone[mem.IO]
// 	if reg&0x00000001 != 0 {

// 	}
// }

func (C *CIA) updateStates() {
	// if C.mem[ICR].IsRead {
	// 	C.mem[ICR].Zone[mem.IO] = 0
	// 	*C.Signal_Pin = 0
	// 	C.mem[ICR].IsRead = false
	// }

	if C.io.Written[ICR] {
		order := C.io.Val[ICR]
		mask := order & 0b00001111
		if mask > 0 {
			if order&0b10000000 > 0 { // 7eme bit = 1 -> mask set
				C.io.CiaRegWrite(ICR, C.io.Val[ICR]|mask)
			} else {
				C.io.CiaRegWrite(ICR, C.io.Val[ICR] & ^mask)
			}
		}
		C.io.Written[ICR] = false
	}

	if C.io.Written[CRA] {
		C.io.Written[CRA] = false
		// Load Latch Once
		if C.io.Val[CRA]&0b00010000 > 0 {
			C.timerAlatch = int32(C.io.Val[TAHI])<<8 + int32(C.io.Val[TALO])
		}
		// Start or stop timer
		if C.io.Val[CRA]&0b00000001 == 1 {
			go C.TimerA()
		}
		if C.io.Val[CRA]&0b00000001 == 0 {
			C.timerAcom <- 1
		}
		C.io.CiaRegWrite(CRA, C.io.Val[CRA]&0b11101111)
	}

	if C.io.Written[CRB] {
		C.io.Written[CRB] = false
		// Load Latch Once
		if C.io.Val[CRB]&0b00010000 > 0 {
			C.timerBlatch = int32(C.io.Val[TBHI])<<8 + int32(C.io.Val[TBLO])
		}
		// Start or stop timer
		if C.io.Val[CRB]&0b00000001 == 1 {
			go C.TimerB()
		}
		if C.io.Val[CRB]&0b00000001 == 0 {
			C.timerBcom <- 1
		}
		C.io.CiaRegWrite(CRB, C.io.Val[CRB]&0b11101111)
	}

	if C.io.Written[TALO] {
		C.io.Written[TALO] = false
		C.timerAlatch = int32(C.io.Val[TAHI])<<8 + int32(C.io.Val[TALO])
	}
	if C.io.Written[TAHI] {
		C.io.Written[TAHI] = false
		C.timerAlatch = int32(C.io.Val[TAHI])<<8 + int32(C.io.Val[TALO])
	}
	if C.io.Written[TBLO] {
		C.io.Written[TBLO] = false
		C.timerBlatch = int32(C.io.Val[TBHI])<<8 + int32(C.io.Val[TBLO])
	}
	if C.io.Written[TBHI] {
		C.io.Written[TBHI] = false
		C.timerBlatch = int32(C.io.Val[TBHI])<<8 + int32(C.io.Val[TBLO])
	}
}

func (C *CIA) Run() {
	C.updateStates()
}

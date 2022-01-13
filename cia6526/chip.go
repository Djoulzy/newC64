package cia6526

import (
	"newC64/memory"
)

type CIA struct {
	name        string
	Reg         [16]byte
	Signal_Pin  *int
	systemCycle *uint16

	timerA_latchLO uint16
	timerA_latchHI uint16
	timerAstate    bool

	timerB_latchLO uint16
	timerB_latchHI uint16
	timerBstate    bool
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
	C.systemCycle = timer

	C.Reg[PRA] = 0x81
	C.Reg[PRB] = 0xFF
	C.Reg[DDRA] = 0x00
	C.Reg[DDRB] = 0x00
	C.Reg[TALO] = 0xFF
	C.Reg[TAHI] = 0xFF
	C.Reg[TBLO] = 0xFF
	C.Reg[TBHI] = 0xFF
	C.Reg[TOD10THS] = 0x00
	C.Reg[TODSEC] = 0x00
	C.Reg[TODMIN] = 0x00
	C.Reg[TODHR] = 0x01
	C.Reg[SRD] = 0x00
	C.Reg[ICR] = 0x00
	C.Reg[CRA] = 0x00
	C.Reg[CRB] = 0x00

	C.timerA_latchLO = 0xFF
	C.timerA_latchHI = 0xFF
	C.timerAstate = false

	C.timerB_latchLO = 0xFF
	C.timerB_latchHI = 0xFF
	C.timerBstate = false
}

func (C *CIA) Run() {
	if C.timerAstate {
		C.TimerA()
	}
	if C.timerBstate {
		C.TimerB()
	}
}

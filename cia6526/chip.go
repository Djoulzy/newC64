package cia6526

import (
	"fmt"
	"newC64/graphic"
	"newC64/memory"
)

type CIA struct {
	name        string
	Reg         [16]byte
	Signal_Pin  *int
	systemCycle *uint16
	InputLine   graphic.KEYPressed
	buffer      Keyboard

	internalPRA    byte
	internalPRB    byte
	timerA_latchLO byte
	timerA_latchHI byte

	timerB_latchLO byte
	timerB_latchHI byte

	interrupt_mask byte

	cycleDelay    int
	VICBankSelect *int
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

const (
	INT_UNDERFL_TA      = 0b00000001
	INT_UNDERFL_TB      = 0b00000010
	INT_ALARM           = 0b00000100
	INT_SDR             = 0b00001000
	INT_INCOMING_SIGNAL = 0b00010000
	INT_SET             = 0b10000000
)

func (C *CIA) Init(name string, memCells *memory.MEM, timer *uint16) {
	C.name = name
	C.systemCycle = timer

	C.Reg[PRA] = 0x00 // 0x81
	C.Reg[PRB] = 0x00 // 0xFF
	C.Reg[DDRA] = 0xFF
	C.Reg[DDRB] = 0xFF
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
	// C.Reg[CRA] = 0x00
	// C.Reg[CRB] = 0x00

	C.timerA_latchLO = 0xFF
	C.timerA_latchHI = 0xFF

	C.timerB_latchLO = 0xFF
	C.timerB_latchHI = 0xFF

	C.interrupt_mask = 0
	C.cycleDelay = 0
}

func (C *CIA) Run() {
	// if C.name == "CIA1" {
	// 	C.buffer = keyMap[charbuff]
	// }
	if C.Reg[CRA]&CTRL_START_STOP > 0 {
		if C.cycleDelay > 0 {
			C.cycleDelay--
			return
		}
		C.TimerA()
	}
	if C.Reg[CRB]&CTRL_START_STOP > 0 {
		C.TimerB()
	}
}

func (C *CIA) Stats() {
	fmt.Printf("%s\n", C.name)
	fmt.Printf("Timer A: ")
	fmt.Printf("Counter: %d - ", (uint16(C.Reg[TAHI])<<8)+uint16(C.Reg[TALO]))
	if C.Reg[CRA]&CTRL_START_STOP > 0 {
		fmt.Printf("Running")
	} else {
		fmt.Printf("Stopped")
	}
	if C.Reg[CRA]&CTRL_LOOP > 0 {
		fmt.Printf(" (Loop)")
	} else {
		fmt.Printf(" (Once)")
	}
	fmt.Printf("\nICR: %04X - CTRL_A: %04X - CTRL_B: %04X", C.Reg[ICR], C.Reg[CRA], C.Reg[CRB])
	fmt.Printf("\nIRQ Line: ")
	if C.Reg[ICR]&INT_SET > 0 {
		fmt.Printf("On ")
	} else {
		fmt.Printf("Off ")
	}
	if C.Reg[ICR]&INT_UNDERFL_TA > 0 {
		fmt.Printf("INT_UNDERFL_A")
	}
	fmt.Println()
}

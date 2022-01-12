package cia6526

import (
	"newC64/memory"
	"newC64/register"
)

type CIA struct {
	name        string
	reg         []register.REG
	Signal_Pin  *int
	systemCycle *uint16

	timerAlatch int32
	timerAstate bool

	timerBlatch int32
	timerBstate bool
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
	C.reg = make([]register.REG, CRB+1)
	C.systemCycle = timer

	C.reg[PRA].Init(memCells, PRA, 0x81)
	C.reg[PRB].Init(memCells, PRB, 0xFF)
	C.reg[DDRA].Init(memCells, DDRA, 0x00)
	C.reg[DDRB].Init(memCells, DDRB, 0x00)
	C.reg[TALO].Init(memCells, TALO, 0xFF)
	C.reg[TAHI].Init(memCells, TAHI, 0xFF)
	C.reg[TBLO].Init(memCells, TBLO, 0xFF)
	C.reg[TBHI].Init(memCells, TBHI, 0xFF)
	C.reg[TOD10THS].Init(memCells, TOD10THS, 0x00)
	C.reg[TODSEC].Init(memCells, TODSEC, 0x00)
	C.reg[TODMIN].Init(memCells, TODMIN, 0x00)
	C.reg[TODHR].Init(memCells, TODHR, 0x01)
	C.reg[SRD].Init(memCells, SRD, 0x00)
	C.reg[ICR].Init(memCells, ICR, 0x00)
	C.reg[CRA].Init(memCells, CRA, 0x00)
	C.reg[CRB].Init(memCells, CRB, 0x00)

	C.timerAlatch = 0
	C.timerBlatch = 0
}

func (C *CIA) updateStates() {
	// if C.mem[ICR].IsRead {
	// 	C.mem[ICR].Zone[mem.IO] = 0
	// 	*C.Signal_Pin = 0
	// 	C.mem[ICR].IsRead = false
	// }

	if C.reg[ICR].IsMofied() {
		order := C.reg[ICR].Input()
		mask := order & 0b00001111
		if mask > 0 {
			if order&0b10000000 > 0 { // 7eme bit = 1 -> mask set
				C.reg[ICR].Output(C.reg[ICR].Val | mask)
			} else {
				C.reg[ICR].Output(C.reg[ICR].Val & ^mask)
			}
		} else {
			C.reg[ICR].Reset()
		}
	}

	if C.reg[CRA].IsMofied() {
		order := C.reg[CRA].Input()
		// Load Latch Once
		if order&0b00010000 > 0 {
			C.timerAlatch = int32(C.reg[TAHI].Val)<<8 + int32(C.reg[TALO].Val)

		}
		// Start or stop timer
		if order&0b00000001 == 1 {
			C.timerAstate = true
		} else {
			C.timerAstate = false
		}
		C.reg[CRA].Output(order & 0b11101111)
	}

	if C.reg[CRB].IsMofied() {
		order := C.reg[CRB].Input()
		// Load Latch Once
		if order&0b00010000 > 0 {
			C.timerAlatch = int32(C.reg[TBHI].Val)<<8 + int32(C.reg[TBLO].Val)
		}
		// Start or stop timer
		if order&0b00000001 == 1 {
			C.timerBstate = true
		} else {
			C.timerBstate = false
		}
		C.reg[CRB].Output(order & 0b11101111)
	}

	if C.reg[TALO].IsMofied() {
		order := C.reg[TALO].Input()
		C.timerAlatch = int32(C.reg[TAHI].Val)<<8 + int32(order)
		C.reg[TALO].Reset()
	}
	if C.reg[TAHI].IsMofied() {
		order := C.reg[TAHI].Input()
		C.timerAlatch = int32(order)<<8 + int32(C.reg[TALO].Val)
		C.reg[TAHI].Reset()
	}
	if C.reg[TBLO].IsMofied() {
		order := C.reg[TBLO].Input()
		C.timerBlatch = int32(C.reg[TBHI].Val)<<8 + int32(order)
		C.reg[TBLO].Reset()
	}
	if C.reg[TBHI].IsMofied() {
		order := C.reg[TBHI].Input()
		C.timerBlatch = int32(order)<<8 + int32(C.reg[TBLO].Val)
		C.reg[TBHI].Reset()
	}
}

func (C *CIA) Run() {
	C.updateStates()
	if C.timerAstate {
		C.TimerA()
	}
	if C.timerBstate {
		C.TimerB()
	}
}

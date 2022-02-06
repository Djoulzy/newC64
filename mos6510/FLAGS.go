package mos6510

import (
	"log"
)

func (C *CPU) bit() {
	var val, oper byte

	switch C.Inst.addr {
	case zeropage:
		oper = C.ram.Read(C.oper)
		val = C.A & oper
		C.setV(oper&0b01000000 == 0b01000000)
		C.setN(oper&0b10000000 == 0b10000000)
	case absolute:
		oper = C.ram.Read(C.oper)
		val = C.A & oper
		C.setV(oper&0b01000000 == 0b01000000)
		C.setN(oper&0b10000000 == 0b10000000)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateZ(val)

}

func (C *CPU) clc() {
	switch C.Inst.addr {
	case implied:
		C.setC(false)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) cld() {
	switch C.Inst.addr {
	case implied:
		C.setD(false)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) cli() {
	switch C.Inst.addr {
	case implied:
		C.setI(false)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) clv() {
	switch C.Inst.addr {
	case implied:
		C.setV(false)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) sec() {
	switch C.Inst.addr {
	case implied:
		C.setC(true)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) sed() {
	// log.Fatal("Decimal mode")
	switch C.Inst.addr {
	case implied:
		C.setD(true)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) sei() {
	switch C.Inst.addr {
	case implied:
		C.setI(true)
	default:
		log.Fatal("Bad addressing mode")
	}

}

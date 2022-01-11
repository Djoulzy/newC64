package mos6510

import (
	"log"
)

func (C *CPU) bit() {
	var val, oper byte

	switch C.inst.addr {
	case zeropage:
		oper = C.ram.Read(C.oper)
		C.setN(oper&^N_mask > 0)
		C.setV(oper&^V_mask > 0)
		val = C.A & oper
	case absolute:
		oper = C.ram.Read(C.oper)
		C.setN(oper&^N_mask > 0)
		C.setV(oper&^V_mask > 0)
		val = C.A & oper
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateZ(val)

}

func (C *CPU) clc() {
	switch C.inst.addr {
	case implied:
		C.setC(false)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) cld() {
	switch C.inst.addr {
	case implied:
		C.setD(false)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) cli() {
	switch C.inst.addr {
	case implied:
		C.setI(false)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) clv() {
	switch C.inst.addr {
	case implied:
		C.setV(false)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) sec() {
	switch C.inst.addr {
	case implied:
		C.setC(true)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) sed() {
	switch C.inst.addr {
	case implied:
		C.setD(true)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) sei() {
	switch C.inst.addr {
	case implied:
		C.setI(true)
	default:
		log.Fatal("Bad addressing mode")
	}

}

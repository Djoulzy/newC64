package mos6510

import (
	"log"
)

func (C *CPU) cmp() {
	var val int

	switch C.inst.addr {
	case immediate:
		val = int(C.A) - int(C.oper)
	case zeropage:
		val = int(C.A) - int(C.ram.Read(C.oper))
	case zeropageX:
		val = int(C.A) - int(C.ram.Read(C.oper+uint16(C.X)))
	case absolute:
		val = int(C.A) - int(C.ram.Read(C.oper))
	case absoluteX:
		val = int(C.A) - int(C.ram.Read(C.oper+uint16(C.X)))
	case absoluteY:
		val = int(C.A) - int(C.ram.Read(C.oper+uint16(C.Y)))
	case indirectX:
		val = int(C.A) - int(C.ReadIndirectX(C.oper))
	case indirectY:
		val = int(C.A) - int(C.ReadIndirectY(C.oper))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setC(val >= 0)
	C.updateN(byte(val))
	C.updateZ(byte(val))

}

func (C *CPU) cpx() {
	var val int

	switch C.inst.addr {
	case immediate:
		val = int(C.X) - int(C.oper)
	case zeropage:
		val = int(C.X) - int(C.ram.Read(C.oper))
	case absolute:
		val = int(C.X) - int(C.ram.Read(C.oper))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setC(val >= 0)
	C.updateN(byte(val))
	C.updateZ(byte(val))

}

func (C *CPU) cpy() {
	var val int

	switch C.inst.addr {
	case immediate:
		val = int(C.Y) - int(C.oper)
	case zeropage:
		val = int(C.Y) - int(C.ram.Read(C.oper))
	case absolute:
		val = int(C.Y) - int(C.ram.Read(C.oper))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setC(val >= 0)
	C.updateN(byte(val))
	C.updateZ(byte(val))

}

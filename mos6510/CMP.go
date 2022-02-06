package mos6510

import (
	"log"
)

func (C *CPU) cmp() {
	var val int
	var crossed bool

	switch C.Inst.addr {
	case immediate:
		val = int(C.A) - int(C.oper)
	case zeropage:
		val = int(C.A) - int(C.ram.Read(C.oper))
	case zeropageX:
		val = int(C.A) - int(C.ram.Read(C.oper+uint16(C.X)))
	case absolute:
		val = int(C.A) - int(C.ram.Read(C.oper))
	case absoluteX:
		C.cross_oper = C.oper + uint16(C.X)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			val = int(C.A) - int(C.ram.Read(C.cross_oper))
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case absoluteY:
		C.cross_oper = C.oper + uint16(C.Y)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			val = int(C.A) - int(C.ram.Read(C.cross_oper))
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case indirectX:
		val = int(C.A) - int(C.ReadIndirectX(C.oper))
	case indirectY:
		C.cross_oper = C.GetIndirectYAddr(C.oper, &crossed)
		if crossed {
			val = int(C.A) - int(C.ram.Read(C.cross_oper))
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		val = int(C.A) - int(C.ram.Read(C.cross_oper))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setC(val >= 0)
	C.updateN(byte(val))
	C.updateZ(byte(val))
}

func (C *CPU) cpx() {
	var val int

	switch C.Inst.addr {
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

	switch C.Inst.addr {
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

package mos6510

import (
	"log"
)

func (C *CPU) lda() {
	var crossed bool

	switch C.Inst.addr {
	case immediate:
		C.A = byte(C.oper)
	case zeropageX:
		C.A = C.ram.Read(C.oper + uint16(C.X))
	case zeropage:
		fallthrough
	case absolute:
		C.A = C.ram.Read(C.oper)
	case absoluteX:
		C.cross_oper = C.oper + uint16(C.X)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			C.A = C.ram.Read(C.cross_oper)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case absoluteY:
		C.cross_oper = C.oper + uint16(C.Y)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			C.A = C.ram.Read(C.cross_oper)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case indirectX:
		C.A = C.ReadIndirectX(C.oper)
	case indirectY:
		C.cross_oper = C.GetIndirectYAddr(C.oper, &crossed)
		if crossed {
			C.A = C.ram.Read(C.cross_oper)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
			return
		}
	case CrossPage:
		C.A = C.ram.Read(C.cross_oper)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)
}

func (C *CPU) sta() {
	switch C.Inst.addr {
	case zeropage:
		C.ram.Write(C.oper, C.A)
	case zeropageX:
		C.ram.Write(C.oper+uint16(C.X), C.A)
	case absolute:
		C.ram.Write(C.oper, C.A)
	case absoluteX:
		C.ram.Write(C.oper+uint16(C.X), C.A)
	case absoluteY:
		C.ram.Write(C.oper+uint16(C.Y), C.A)
	case indirectX:
		C.WriteIndirectX(C.oper, C.A)
	case indirectY:
		C.WriteIndirectY(C.oper, C.A)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) ldx() {
	switch C.Inst.addr {
	case immediate:
		C.X = byte(C.oper)
	case zeropage:
		C.X = C.ram.Read(C.oper)
	case zeropageY:
		C.X = C.ram.Read(C.oper + uint16(C.Y))
	case absolute:
		C.X = C.ram.Read(C.oper)
	case absoluteY:
		C.X = C.ram.Read(C.oper + uint16(C.Y))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.X)
	C.updateZ(C.X)

}

func (C *CPU) stx() {
	switch C.Inst.addr {
	case zeropage:
		C.ram.Write(C.oper, C.X)
	case zeropageY:
		C.ram.Write(C.oper+uint16(C.Y), C.X)
	case absolute:
		C.ram.Write(C.oper, C.X)
	default:
		log.Fatal("Bad addressing mode")
	}

}

func (C *CPU) ldy() {
	switch C.Inst.addr {
	case immediate:
		C.Y = byte(C.oper)
	case zeropage:
		C.Y = C.ram.Read(C.oper)
	case zeropageX:
		C.Y = C.ram.Read(C.oper + uint16(C.X))
	case absolute:
		C.Y = C.ram.Read(C.oper)
	case absoluteX:
		C.Y = C.ram.Read(C.oper + uint16(C.X))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.Y)
	C.updateZ(C.Y)

}

func (C *CPU) sty() {
	switch C.Inst.addr {
	case zeropage:
		C.ram.Write(C.oper, C.Y)
	case zeropageX:
		C.ram.Write(C.oper+uint16(C.X), C.Y)
	case absolute:
		C.ram.Write(C.oper, C.Y)
	default:
		log.Fatal("Bad addressing mode")
	}

}

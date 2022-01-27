package mos6510

import (
	"log"
)

func (C *CPU) adc() {
	var val uint16
	var oper byte

	// log.Printf("%04X - %s", C.InstStart, C.registers())
	// log.Fatal("ADC")
	switch C.Inst.addr {
	case immediate:
		val = uint16(C.A) + C.oper + uint16(C.getC())
		C.setC(val > 0x00FF)
		C.updateV(C.A, byte(oper), byte(val))
		C.A = byte(val)
	case zeropage:
		fallthrough
	case absolute:
		oper = C.ram.Read(C.oper)
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case zeropageX:
		oper = C.ram.Read(C.oper + uint16(C.X))
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case absoluteX:
		C.cross_oper = C.oper + uint16(C.X)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			oper = C.ram.Read(C.cross_oper)
			val = uint16(C.A) + uint16(oper) + uint16(C.getC())
			C.setC(val > 0x00FF)
			C.updateV(C.A, oper, byte(val))
			C.A = byte(val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
		}
	case absoluteY:
		C.cross_oper = C.oper + uint16(C.Y)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			oper = C.ram.Read(C.cross_oper)
			val = uint16(C.A) + uint16(oper) + uint16(C.getC())
			C.setC(val > 0x00FF)
			C.updateV(C.A, oper, byte(val))
			C.A = byte(val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
		}
	case indirectX:
		oper = C.ReadIndirectX(C.oper)
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case indirectY:
		C.cross_oper = C.GetIndirectYAddr(C.oper)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			oper = C.ram.Read(C.cross_oper)
			val = uint16(C.A) + uint16(oper) + uint16(C.getC())
			C.setC(val > 0x00FF)
			C.updateV(C.A, oper, byte(val))
			C.A = byte(val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
		}
	case CrossPage:
		oper = C.ram.Read(C.cross_oper)
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(byte(val))
	C.updateZ(byte(val))

}

func (C *CPU) sbc() {
	var val int
	var oper byte

	switch C.Inst.addr {
	case immediate:
		val = int(C.A) - int(C.oper)
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^byte(C.oper), byte(val))
		C.A = byte(val)
	case zeropage:
		fallthrough
	case absolute:
		val = int(C.A) - int(C.ram.Read(C.oper))
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^C.ram.Read(C.oper), byte(val))
		C.A = byte(val)
	case zeropageX:
		fallthrough
	case absoluteX:
		val = int(C.A) - int(C.ram.Read(C.oper+uint16(C.X)))
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^C.ram.Read(C.oper+uint16(C.X)), byte(val))
		C.A = byte(val)
	case absoluteY:
		val = int(C.A) - int(C.ram.Read(C.oper+uint16(C.Y)))
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^C.ram.Read(C.oper+uint16(C.Y)), byte(val))
		C.A = byte(val)
	case indirectX:
		val = int(C.A) - int(C.ReadIndirectX(C.oper))
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^C.ReadIndirectX(C.oper), byte(val))
		C.A = byte(val)
	case indirectY:
		C.cross_oper = C.GetIndirectYAddr(C.oper)
		oper = C.ram.Read(C.cross_oper)
		if C.oper&0xFF00 == C.cross_oper&0xFF00 {
			val = int(C.A) - int(oper)
			if C.getC() == 0 {
				val -= 1
			}
			C.updateV(C.A, ^oper, byte(val))
			C.A = byte(val)
		} else {
			C.Inst.addr = CrossPage
			C.State = Compute
			C.Inst.Cycles++
		}
	case CrossPage:
		oper = C.ram.Read(C.cross_oper)
		val = int(C.A) - int(oper)
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^oper, byte(val))
		C.A = byte(val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setC(val >= 0x00)
	C.setN(val&0b10000000 == 0b10000000)
}

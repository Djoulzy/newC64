package mos6510

import (
	"log"
)

func (C *CPU) adc() {
	var val uint16
	var oper byte

	// log.Printf("%04X - %s", C.InstStart, C.registers())
	// log.Fatal("ADC")
	switch C.inst.addr {
	case immediate:
		val = uint16(C.A) + C.oper + uint16(C.getC())
		C.setC(val > 0x00FF)
		C.updateV(C.A, byte(oper), byte(val))
		C.A = byte(val)
	case zeropage:
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
	case absolute:
		oper = C.ram.Read(C.oper)
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case absoluteX:
		oper = C.ram.Read(C.oper + uint16(C.X))
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case absoluteY:
		oper = C.ram.Read(C.oper + uint16(C.Y))
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case indirectX:
		oper = C.ReadIndirectX(C.oper)
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x00FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case indirectY:
		oper = C.ReadIndirectY(C.oper)
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

	switch C.inst.addr {
	case immediate:
		val = int(C.A) - int(C.oper)
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^byte(C.oper), byte(val))
		C.A = byte(val)
	case zeropage:
		val = int(C.A) - int(C.ram.Read(C.oper))
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^C.ram.Read(C.oper), byte(val))
		C.A = byte(val)
	case zeropageX:
		val = int(C.A) - int(C.ram.Read(C.oper+uint16(C.X)))
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^C.ram.Read(C.oper+uint16(C.X)), byte(val))
		C.A = byte(val)
	case absolute:
		val = int(C.A) - int(C.ram.Read(C.oper))
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^C.ram.Read(C.oper), byte(val))
		C.A = byte(val)
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
		val = int(C.A) - int(C.ReadIndirectY(C.oper))
		if C.getC() == 0 {
			val -= 1
		}
		C.updateV(C.A, ^C.ReadIndirectY(C.oper), byte(val))
		C.A = byte(val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setC(val >= 0x00)
	C.setN(val&0b10000000 == 0b10000000)
}

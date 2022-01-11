package mos6510

import (
	"fmt"
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
	if C.conf.Disassamble {
		fmt.Printf("\n")
	}
}

func (C *CPU) sbc() {
	var val uint16
	var oper byte

	switch C.inst.addr {
	case immediate:
		val = uint16(C.A) + ^C.oper + uint16(C.getC())
		C.setC(val > 0x0FF)
		C.updateV(C.A, byte(^oper), byte(val))
		C.A = byte(val)
	case zeropage:
		oper = ^C.ram.Read(C.oper)
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x0FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case zeropageX:
		oper = ^C.ram.Read(C.oper + uint16(C.X))
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x0FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case absolute:
		oper = ^C.ram.Read(C.oper)
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x0FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case absoluteX:
		oper = ^C.ram.Read(C.oper + uint16(C.X))
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x0FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case absoluteY:
		oper = ^C.ram.Read(C.oper + uint16(C.Y))
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x0FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case indirectX:
		oper = ^C.ReadIndirectX(C.oper)
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x0FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	case indirectY:
		oper = ^C.ReadIndirectY(C.oper)
		val = uint16(C.A) + uint16(oper) + uint16(C.getC())
		C.setC(val > 0x0FF)
		C.updateV(C.A, oper, byte(val))
		C.A = byte(val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(byte(val))
	C.updateZ(byte(val))
	if C.conf.Disassamble {
		fmt.Printf("\n")
	}
}

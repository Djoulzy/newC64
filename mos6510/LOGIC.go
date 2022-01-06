package mos6510

import (
	"fmt"
	"log"
)

func (C *CPU) and() {
	switch C.inst.addr {
	case immediate:
		C.A &= byte(C.oper)
	case zeropage:
		C.A &= C.ram.Read(C.oper)
	case zeropageX:
		C.A &= C.ram.Read(C.oper + uint16(C.X))
	case absolute:
		C.A &= C.ram.Read(C.oper)
	case absoluteX:
		C.A &= C.ram.Read(C.oper + uint16(C.X))
	case absoluteY:
		C.A &= C.ram.Read(C.oper + uint16(C.Y))
	case indirectX:
		C.A &= C.ReadIndirectX(C.oper)
	case indirectY:
		C.A &= C.ReadIndirectY(C.oper)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)
	fmt.Printf("\n")
}

func (C *CPU) asl() {
	var val uint16

	switch C.inst.addr {
	case implied:
		val = uint16(C.A) << 1
		C.A = byte(val)
	case zeropage:
		val = uint16(C.ram.Read(C.oper)) << 1
		C.ram.Write(C.oper, byte(val))
	case zeropageX:
		val = uint16(C.ram.Read(C.oper+uint16(C.X))) << 1
		C.ram.Write(C.oper, byte(val))
	case absolute:
		val = uint16(C.ram.Read(C.oper)) << 1
		C.ram.Write(C.oper, byte(val))
	case absoluteX:
		val = uint16(C.ram.Read(C.oper+uint16(C.X))) << 1
		C.ram.Write(C.oper, byte(val))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(byte(val))
	C.updateZ(byte(val))
	C.setC(val > 0x00FF)
	fmt.Printf("\n")
}

func (C *CPU) eor() {
	fmt.Printf("Not implemented: %v\n", C.inst)
}

func (C *CPU) lsr() {
	var val byte

	switch C.inst.addr {
	case implied:
		C.setC(C.A&0x01 == 0x01)
		val = C.A >> 1
		C.A = val
	case zeropage:
		val = C.ram.Read(C.oper)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		C.ram.Write(C.oper, val)
	case zeropageX:
		val = C.ram.Read(C.oper + uint16(C.X))
		C.setC(val&0x01 == 0x01)
		val >>= 1
		C.ram.Write(C.oper, val)
	case absolute:
		val = C.ram.Read(C.oper)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		C.ram.Write(C.oper, val)
	case absoluteX:
		val = C.ram.Read(C.oper + uint16(C.X))
		C.setC(val&0x01 == 0x01)
		val >>= 1
		C.ram.Write(C.oper, val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setN(false)
	C.updateZ(byte(val))
	fmt.Printf("\n")
}

func (C *CPU) ora() {
	fmt.Printf("Not implemented: %v\n", C.inst)
}

func (C *CPU) rla() {
	fmt.Printf("Not implemented: %v\n", C.inst)
}

func (C *CPU) rol() {
	var val uint16

	switch C.inst.addr {
	case implied:
		val = uint16(C.A) << 1
		if C.issetC() {
			val++
		}
		C.A = byte(val)
	case zeropage:
		val = uint16(C.ram.Read(C.oper)) << 1
		if C.issetC() {
			val++
		}
		C.ram.Write(C.oper, byte(val))
	case zeropageX:
		val = uint16(C.ram.Read(C.oper+uint16(C.X))) << 1
		if C.issetC() {
			val++
		}
		C.ram.Write(C.oper, byte(val))
	case absolute:
		val = uint16(C.ram.Read(C.oper)) << 1
		if C.issetC() {
			val++
		}
		C.ram.Write(C.oper, byte(val))
	case absoluteX:
		val = uint16(C.ram.Read(C.oper+uint16(C.X))) << 1
		if C.issetC() {
			val++
		}
		C.ram.Write(C.oper, byte(val))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(byte(val))
	C.updateZ(byte(val))
	C.setC(val > 0x00FF)
	fmt.Printf("\n")
}

func (C *CPU) ror() {
	var val byte

	switch C.inst.addr {
	case implied:
		C.setC(C.A&0x01 == 0x01)
		val = C.A >> 1
		if C.issetC() {
			C.A |= 0b10000000
		}
		C.A = val
	case zeropage:
		val = C.ram.Read(C.oper)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.ram.Write(C.oper, val)
	case zeropageX:
		val = C.ram.Read(C.oper + uint16(C.X))
		C.setC(val&0x01 == 0x01)
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.ram.Write(C.oper, val>>1)
	case absolute:
		val = C.ram.Read(C.oper)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.ram.Write(C.oper, val)
	case absoluteX:
		val = C.ram.Read(C.oper + uint16(C.X))
		C.setC(val&0x01 == 0x01)
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.ram.Write(C.oper, val>>1)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setN(false)
	C.updateZ(byte(val))
	fmt.Printf("\n")
}

func (C *CPU) sax() {
	fmt.Printf("Not implemented: %v\n", C.inst)
}

func (C *CPU) slo() {
	fmt.Printf("Not implemented: %v\n", C.inst)
}

func (C *CPU) sre() {
	fmt.Printf("Not implemented: %v\n", C.inst)
}

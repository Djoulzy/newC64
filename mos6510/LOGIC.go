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
		dest := C.oper+uint16(C.X)
		val = uint16(C.ram.Read(dest)) << 1
		C.ram.Write(dest, byte(val))
	case absolute:
		val = uint16(C.ram.Read(C.oper)) << 1
		C.ram.Write(C.oper, byte(val))
	case absoluteX:
		dest := C.oper+uint16(C.X)
		val = uint16(C.ram.Read(dest)) << 1
		C.ram.Write(dest, byte(val))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(byte(val))
	C.updateZ(byte(val))
	C.setC(val > 0x00FF)

}

func (C *CPU) eor() {
	switch C.inst.addr {
	case immediate:
		C.A ^= byte(C.oper)
	case zeropage:
		C.A ^= C.ram.Read(C.oper)
	case zeropageX:
		C.A ^= C.ram.Read(C.oper + uint16(C.X))
	case absolute:
		C.A ^= C.ram.Read(C.oper)
	case absoluteX:
		C.A ^= C.ram.Read(C.oper + uint16(C.X))
	case absoluteY:
		C.A ^= C.ram.Read(C.oper + uint16(C.Y))
	case indirectX:
		C.A ^= C.ReadIndirectX(C.oper)
	case indirectY:
		C.A ^= C.ReadIndirectY(C.oper)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)

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
		dest := C.oper + uint16(C.X)
		val = C.ram.Read(dest)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		C.ram.Write(dest, val)
	case absolute:
		val = C.ram.Read(C.oper)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		C.ram.Write(C.oper, val)
	case absoluteX:
		dest := C.oper + uint16(C.X)
		val = C.ram.Read(dest)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		C.ram.Write(dest, val)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setN(false)
	C.updateZ(byte(val))

}

func (C *CPU) ora() {
	switch C.inst.addr {
	case immediate:
		C.A |= byte(C.oper)
	case zeropage:
		C.A |= C.ram.Read(C.oper)
	case zeropageX:
		C.A |= C.ram.Read(C.oper + uint16(C.X))
	case absolute:
		C.A |= C.ram.Read(C.oper)
	case absoluteX:
		C.A |= C.ram.Read(C.oper + uint16(C.X))
	case absoluteY:
		C.A |= C.ram.Read(C.oper + uint16(C.Y))
	case indirectX:
		C.A |= C.ReadIndirectX(C.oper)
	case indirectY:
		C.A |= C.ReadIndirectY(C.oper)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)

}

func (C *CPU) rla() {
	fmt.Printf("%s\nNot implemented: %v\n", C.Disassemble(), C.inst)
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
		dest := C.oper+uint16(C.X)
		val = uint16(C.ram.Read(dest)) << 1
		if C.issetC() {
			val++
		}
		C.ram.Write(dest, byte(val))
	case absolute:
		val = uint16(C.ram.Read(C.oper)) << 1
		if C.issetC() {
			val++
		}
		C.ram.Write(C.oper, byte(val))
	case absoluteX:
		dest := C.oper+uint16(C.X)
		val = uint16(C.ram.Read(dest)) << 1
		if C.issetC() {
			val++
		}
		C.ram.Write(dest, byte(val))
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(byte(val))
	C.updateZ(byte(val))
	C.setC(val > 0x00FF)

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
		dest := C.oper + uint16(C.X)
		val = C.ram.Read(dest)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.ram.Write(dest, val>>1)
	case absolute:
		val = C.ram.Read(C.oper)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.ram.Write(C.oper, val)
	case absoluteX:
		dest := C.oper + uint16(C.X)
		val = C.ram.Read(dest)
		C.setC(val&0x01 == 0x01)
		val >>= 1
		if C.issetC() {
			val |= 0b10000000
		}
		C.ram.Write(dest, val>>1)
	default:
		log.Fatal("Bad addressing mode")
	}
	C.setN(false)
	C.updateZ(byte(val))

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

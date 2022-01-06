package mos6510

import (
	"fmt"
	"log"
)

func (C *CPU) bcc() {
	dest := C.getRelativeAddr(C.oper)
	switch C.inst.addr {
	case relative:
		if !C.issetC() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("(%04X)\n", dest)
}

func (C *CPU) bcs() {
	dest := C.getRelativeAddr(C.oper)
	switch C.inst.addr {
	case relative:
		if C.issetC() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("(%04X)\n", dest)
}

func (C *CPU) beq() {
	dest := C.getRelativeAddr(C.oper)
	switch C.inst.addr {
	case relative:
		if C.issetZ() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("(%04X)\n", dest)
}

func (C *CPU) bmi() {
	dest := C.getRelativeAddr(C.oper)
	switch C.inst.addr {
	case relative:
		if C.issetN() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("(%04X)\n", dest)
}

func (C *CPU) bne() {
	dest := C.getRelativeAddr(C.oper)
	switch C.inst.addr {
	case relative:
		if !C.issetZ() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("(%04X)\n", dest)
}

func (C *CPU) bpl() {
	dest := C.getRelativeAddr(C.oper)
	switch C.inst.addr {
	case relative:
		if !C.issetN() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("(%04X)\n", dest)
}

func (C *CPU) brk() {
	switch C.inst.addr {
	case implied:
		C.pushWordStack(C.PC + 1)
		C.setB(true)
		C.pushByteStack(C.S)
		C.PC = C.readWord(IRQBRK_Vector)
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}

func (C *CPU) bvc() {
	dest := C.getRelativeAddr(C.oper)
	switch C.inst.addr {
	case relative:
		if !C.issetV() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("(%04X)\n", dest)
}

func (C *CPU) bvs() {
	dest := C.getRelativeAddr(C.oper)
	switch C.inst.addr {
	case relative:
		if C.issetV() {
			C.PC = C.getRelativeAddr(C.oper)
		}
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("(%04X)\n", dest)
}

func (C *CPU) jmp() {
	switch C.inst.addr {
	case absolute:
		C.PC = C.oper
	case indirect:
		C.PC = C.readWord(C.oper)
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}

func (C *CPU) jsr() {
	switch C.inst.addr {
	case absolute:
		C.pushWordStack(C.PC - 1)
		C.PC = C.oper
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}

func (C *CPU) rti() {
	switch C.inst.addr {
	case implied:
		C.S = C.pullByteStack()
		C.setB(false)
		C.setU(false)
		C.PC = C.pullWordStack()
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}

func (C *CPU) rts() {
	switch C.inst.addr {
	case implied:
		C.PC = C.pullWordStack() + 1
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}

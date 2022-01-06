package mos6510

import (
	"fmt"
	"log"
)

func (C *CPU) pha() {
	switch C.inst.addr {
	case implied:
		C.pushByteStack(C.A)
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}

func (C *CPU) php() {
	switch C.inst.addr {
	case implied:
		C.setB(true)
		C.setU(true)
		C.pushByteStack(C.S)
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}

func (C *CPU) pla() {
	switch C.inst.addr {
	case implied:
		C.A = C.pullByteStack()
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)
	fmt.Printf("\n")
}

func (C *CPU) plp() {
	switch C.inst.addr {
	case implied:
		C.S = C.pullByteStack()
		C.setB(false)
		C.setU(false)
	default:
		log.Fatal("Bad addressing mode")
	}
	fmt.Printf("\n")
}

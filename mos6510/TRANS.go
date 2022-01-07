package mos6510

import (
	"fmt"
	"log"
)

func (C *CPU) tax() {
	switch C.inst.addr {
	case implied:
		C.X = C.A
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.X)
	C.updateZ(C.X)
	if C.conf.Disassamble {
		fmt.Printf("\n")
	}
}

func (C *CPU) tay() {
	switch C.inst.addr {
	case implied:
		C.Y = C.A
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.Y)
	C.updateZ(C.Y)
	if C.conf.Disassamble {
		fmt.Printf("\n")
	}
}

func (C *CPU) tsx() {
	switch C.inst.addr {
	case implied:
		C.X = C.SP
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.X)
	C.updateZ(C.X)
	if C.conf.Disassamble {
		fmt.Printf("\n")
	}
}

func (C *CPU) txa() {
	switch C.inst.addr {
	case implied:
		C.A = C.X
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)
	if C.conf.Disassamble {
		fmt.Printf("\n")
	}
}

func (C *CPU) txs() {
	switch C.inst.addr {
	case implied:
		C.SP = C.X
	default:
		log.Fatal("Bad addressing mode")
	}
	if C.conf.Disassamble {
		fmt.Printf("\n")
	}
}

func (C *CPU) tya() {
	switch C.inst.addr {
	case implied:
		C.Y = C.X
	default:
		log.Fatal("Bad addressing mode")
	}
	C.updateN(C.A)
	C.updateZ(C.A)
	if C.conf.Disassamble {
		fmt.Printf("\n")
	}
}

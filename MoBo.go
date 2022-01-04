package main

import (
	"newC64/mos6510"
	"newC64/pla906114"
	"newC64/ram4164"
	"newC64/roms"
)

var (
	cpu     CPU
	mem     MEM
	pla     PLA
	kernal  ROM
	basic   ROM
	chargen ROM
)

func setup() {
	// ROMs & RAM Setup
	mem = &ram4164.RAM{}
	kernal = &roms.ROM{}
	kernal.Init("assets/roms/kernal.bin", 8192)
	basic = &roms.ROM{}
	basic.Init("assets/roms/basic.bin", 8192)
	chargen = &roms.ROM{}
	chargen.Init("assets/roms/char.bin", 4096)

	// PLA Setup
	pla = &pla906114.PLA{}
	pla.Attach(mem, pla906114.RAM)
	pla.Attach(kernal, pla906114.KERNAL)
	pla.Attach(basic, pla906114.BASIC)
	pla.Attach(basic, pla906114.CHAR)
	pla.Init()

	// CPU Setup
	cpu = &mos6510.CPU{}
	cpu.Init(pla)
}

func main() {
	setup()
}

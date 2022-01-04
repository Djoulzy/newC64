package main

import (
	"newC64/memory"
	"newC64/mos6510"
	"newC64/pla906114"
)

const (
	memorySize = 65536
)

var (
	cpu     CPU
	mem     MEM
	pla     PLA
	kernal  MEM
	basic   MEM
	chargen MEM
)

func setup() {
	// ROMs & RAM Setup
	mem = &memory.MEM{Size: memorySize}
	kernal = &memory.MEM{Size: 8192}
	kernal.Load("assets/roms/kernal.bin")
	basic = &memory.MEM{Size: 8192}
	basic.Load("assets/roms/basic.bin")
	chargen = &memory.MEM{Size: 4096}
	chargen.Load("assets/roms/char.bin")

	// PLA Setup
	pla = &pla906114.PLA{}
	pla.Attach(mem, pla906114.RAM)
	pla.Attach(kernal, pla906114.KERNAL)
	pla.Attach(basic, pla906114.BASIC)
	pla.Attach(basic, pla906114.CHAR)

	// CPU Setup
	cpu = &mos6510.CPU{}
	cpu.Init(pla)

	pla.Dump(0x0000)
	pla.Dump(0xE000)
}

func main() {
	setup()
}

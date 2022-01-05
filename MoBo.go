package main

import (
	"newC64/graphic"
	"newC64/memory"
	"newC64/mos6510"
	"newC64/pla906114"
	"newC64/vic6569"
	"runtime"
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
	io      MEM
	vic     vic6569.VIC

	video graphic.Driver
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func setup() {
	// ROMs & RAM Setup
	mem = &memory.MEM{Size: memorySize}
	io = &memory.MEM{Size: 4096}
	kernal = &memory.MEM{Size: 8192}
	kernal.Load("assets/roms/kernal.bin")
	basic = &memory.MEM{Size: 8192}
	basic.Load("assets/roms/basic.bin")
	chargen = &memory.MEM{Size: 4096}
	chargen.Load("assets/roms/char.bin")

	// PLA Setup
	pla = &pla906114.PLA{}
	pla.Attach(mem, pla906114.RAM, 0)
	pla.Attach(io, pla906114.IO, pla906114.IOStart)
	pla.Attach(kernal, pla906114.KERNAL, pla906114.KernalStart)
	pla.Attach(basic, pla906114.BASIC, pla906114.BasicStart)
	pla.Attach(basic, pla906114.CHAR, pla906114.CharStart)

	// CPU Setup
	cpu = &mos6510.CPU{}
	cpu.Init(pla)

	pla.Dump(0x0000)
	pla.Dump(0xFEFF)

	video = &graphic.SDLDriver{}
	vic = vic6569.VIC{}
	vic.Init(mem, io, chargen, video)
}

func main() {
	setup()

	for {
		vic.Run()
		cpu.NextCycle()
	}
}

package main

import (
	"github.com/Djoulzy/Tools/clog"

	"github.com/Djoulzy/emutools/mem"
)

const (
	nbMemLayout = 3

	ramSize     = 65536
	kernalSize  = 8192
	basicSize   = 8192
	ioSize      = 4096
	chargenSize = 4096

	StartLogging = true
	LogLevel     = 5
)

var (
	KernalAccess mem.MEMAccess

	RAM     []byte
	IO      []byte
	KERNAL  []byte
	BASIC   []byte
	CHARGEN []byte

	BankSel byte
	Mem     mem.BANK
)

type accessor struct {
}

func (C *accessor) MRead(mem []byte, addr uint16) byte {
	clog.Test("Accessor", "MRead", "Addr: %04X", addr)
	return mem[addr]
}

func (C *accessor) MWrite(meme []byte, addr uint16, val byte) {

}

func setup() {
	RAM = make([]byte, ramSize)
	IO = make([]byte, ioSize)
	KERNAL = mem.LoadROM(kernalSize, "assets/roms/kernal.bin")
	BASIC = mem.LoadROM(basicSize, "assets/roms/basic.bin")
	CHARGEN = mem.LoadROM(chargenSize, "assets/roms/char.bin")

	Mem = mem.InitBanks(nbMemLayout, &BankSel)
	KernalAccess = &accessor{}
}

func layout() {
	Mem.Layouts[0] = mem.InitConfig(ramSize)
	Mem.Layouts[0].Attach("RAM", 0, RAM, mem.READWRITE, false)
	Mem.Layouts[0].Attach("KERNAL", 14, KERNAL, mem.READONLY, false)
	Mem.Layouts[0].Attach("BASIC", 10, BASIC, mem.READONLY, false)
	Mem.Layouts[0].Attach("CHARGEN", 13, CHARGEN, mem.READONLY, false)
	Mem.Layouts[0].Attach("IO", 4, IO, mem.READWRITE, false)
	Mem.Layouts[0].Accessor("KERNAL", KernalAccess)
	Mem.Layouts[0].Show()

	Mem.Layouts[1] = mem.InitConfig(ramSize)
	Mem.Layouts[1].Attach("RAM", 0, RAM, mem.READWRITE, false)
	Mem.Layouts[1].Attach("KERNAL", 14, KERNAL, mem.READONLY, false)
	Mem.Layouts[1].Attach("CHARGEN", 13, CHARGEN, mem.READONLY, false)
	Mem.Layouts[1].Show()

	Mem.Layouts[2] = mem.InitConfig(ramSize)
	Mem.Layouts[2].Attach("RAM", 0, RAM, mem.READWRITE, false)
	Mem.Layouts[2].Attach("IO", 13, IO, mem.READWRITE, false)
	Mem.Layouts[2].Show()
}

func main() {
	clog.LogLevel = LogLevel
	clog.StartLogging = StartLogging

	setup()
	layout()

	BankSel = 2

	Mem.Write(0xD0FF, 0xEE)
	Mem.Read(0xD0FF)
	Mem.Read(0xD0FF)

	Mem.Dump(0xD0FF)
	Mem.Dump(0xD0FF)
	// fmt.Printf("Read: %04X\n", Mem.Read(0xF000))
	// fmt.Printf("Read: %04X\n", Mem.Read(0x5000))
}
package main

import (
	"newC64/mem"
)

type accessor struct {
}

func (C *accessor) MRead(mem []byte, addr uint16) byte {
	// clog.Test("Accessor", "MRead", "Addr: %04X", addr)
	return IORead_Mapper[addr](addr)
}

func (C *accessor) MWrite(meme []byte, addr uint16, val byte) {
	// clog.Test("Accessor", "MWrite", "Addr: %04X", addr)
	IOWrite_Mapper[addr](addr, val)
}

func memLayouts() {
	MEM.Layouts[31] = mem.InitConfig(4, ramSize)
	MEM.Layouts[31].Attach("RAM", 0, 0, RAM, mem.READWRITE)
	MEM.Layouts[31].Attach("BASIC", 1, 10, BASIC, mem.READONLY)
	MEM.Layouts[31].Attach("IO", 2, 13, IO, mem.READWRITE)
	MEM.Layouts[31].Attach("KERNAL", 3, 14, KERNAL, mem.READONLY)
	MEM.Layouts[31].Accessor(2, IOAccess)
	MEM.Layouts[31].Show()

	MEM.Layouts[26] = mem.InitConfig(3, ramSize)
	MEM.Layouts[26].Attach("RAM", 0, 0, RAM, mem.READWRITE)
	MEM.Layouts[26].Attach("CHARGEN", 1, 13, CHARGEN, mem.READONLY)
	MEM.Layouts[26].Attach("KERNAL", 2, 14, KERNAL, mem.READONLY)
	MEM.Layouts[26].Show()

	MEM.Layouts[29] = mem.InitConfig(2, ramSize)
	MEM.Layouts[29].Attach("RAM", 0, 0, RAM, mem.READWRITE)
	MEM.Layouts[29].Attach("IO", 1, 13, IO, mem.READWRITE)
	MEM.Layouts[29].Accessor(1, IOAccess)
	MEM.Layouts[29].Show()

	// MEM.Layouts[7] = mem.InitConfig(5, ramSize)
	// MEM.Layouts[7].Attach("RAM", 0, 0, RAM, mem.READWRITE)
	// MEM.Layouts[7].Attach("CART_LO", 1, 8, CHARGEN, mem.READONLY)
	// MEM.Layouts[7].Attach("CART_HI", 2, 10, CHARGEN, mem.READONLY)
	// MEM.Layouts[7].Attach("IO", 3, 13, IO, mem.READWRITE)
	// MEM.Layouts[7].Attach("KERNAL", 4, 14, KERNAL, mem.READONLY)
	// MEM.Layouts[7].Show()

	MEM.Layouts[7] = mem.InitConfig(3, ramSize)
	MEM.Layouts[7].Attach("RAM", 0, 0, RAM, mem.READWRITE)
	MEM.Layouts[7].Attach("IO", 1, 13, IO, mem.READWRITE)
	MEM.Layouts[7].Attach("KERNAL", 2, 14, KERNAL, mem.READONLY)
	MEM.Layouts[7].Accessor(1, IOAccess)
	MEM.Layouts[7].Show()

	MEM.Layouts[0] = mem.InitConfig(1, ramSize)
	MEM.Layouts[0].Attach("RAM", 0, 0, RAM, mem.READWRITE)
}

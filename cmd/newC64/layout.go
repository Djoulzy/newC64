package main

import (
	"github.com/Djoulzy/emutools/mem"
)

func memLayouts() {
	MEM.Layouts[31] = mem.InitConfig(ramSize)
	MEM.Layouts[31].Attach("RAM", 0, RAM, mem.READWRITE)
	MEM.Layouts[31].Attach("BASIC", 0xA000, BASIC, mem.READONLY)
	MEM.Layouts[31].Attach("IO", 0xD000, IO, mem.READWRITE)
	MEM.Layouts[31].Attach("KERNAL", 0xE000, KERNAL, mem.READONLY)
	MEM.Layouts[31].Accessor("RAM", RAMAccess)
	MEM.Layouts[31].Accessor("IO", IOAccess)
	// MEM.Layouts[31].Show()

	MEM.Layouts[30] = mem.InitConfig(ramSize)
	MEM.Layouts[30].Attach("RAM", 0, RAM, mem.READWRITE)
	MEM.Layouts[30].Attach("IO", 0xD000, IO, mem.READWRITE)
	MEM.Layouts[30].Attach("KERNAL", 0xE000, KERNAL, mem.READONLY)
	MEM.Layouts[30].Accessor("RAM", RAMAccess)
	MEM.Layouts[30].Accessor("IO", IOAccess)
	// MEM.Layouts[30].Show()

	MEM.Layouts[29] = mem.InitConfig(ramSize)
	MEM.Layouts[29].Attach("RAM", 0, RAM, mem.READWRITE)
	MEM.Layouts[29].Attach("IO", 0xD000, IO, mem.READWRITE)
	MEM.Layouts[29].Accessor("RAM", RAMAccess)
	MEM.Layouts[29].Accessor("IO", IOAccess)
	// MEM.Layouts[29].Show()

	MEM.Layouts[28] = mem.InitConfig(ramSize)
	MEM.Layouts[28].Attach("RAM", 0, RAM, mem.READWRITE)
	MEM.Layouts[28].Accessor("RAM", RAMAccess)

	MEM.Layouts[26] = mem.InitConfig(ramSize)
	MEM.Layouts[26].Attach("RAM", 0, RAM, mem.READWRITE)
	MEM.Layouts[26].Attach("CHARGEN", 0xD000, CHARGEN, mem.READONLY)
	MEM.Layouts[26].Attach("KERNAL", 0xE000, KERNAL, mem.READONLY)
	MEM.Layouts[26].Accessor("RAM", RAMAccess)
	// MEM.Layouts[26].Show()

	MEM.Layouts[25] = mem.InitConfig(ramSize)
	MEM.Layouts[25].Attach("RAM", 0, RAM, mem.READWRITE)
	MEM.Layouts[25].Attach("CHARGEN", 0xD000, CHARGEN, mem.READONLY)
	MEM.Layouts[25].Accessor("RAM", RAMAccess)
	// MEM.Layouts[25].Show()

	// MEM.Layouts[7] = mem.InitConfig(5, ramSize)
	// MEM.Layouts[7].Attach("RAM", 0, 0, RAM, mem.READWRITE)
	// MEM.Layouts[7].Attach("CART_LO", 1, 8, CHARGEN, mem.READONLY)
	// MEM.Layouts[7].Attach("CART_HI", 2, 10, CHARGEN, mem.READONLY)
	// MEM.Layouts[7].Attach("IO", 3, 13, IO, mem.READWRITE)
	// MEM.Layouts[7].Attach("KERNAL", 4, 14, KERNAL, mem.READONLY)
	// MEM.Layouts[7].Show()

	// MEM.Layouts[7] = mem.InitConfig(4, ramSize)
	// MEM.Layouts[7].Attach("RAM", 0, 0, RAM, mem.READWRITE)
	// MEM.Layouts[7].Attach("BASIC", 1, 10, BASIC, mem.READONLY)
	// MEM.Layouts[7].Attach("IO", 2, 13, IO, mem.READWRITE)
	// MEM.Layouts[7].Attach("KERNAL", 3, 14, KERNAL, mem.READONLY)
	// MEM.Layouts[7].Accessor(2, IOAccess)
	// MEM.Layouts[7].Show()

	MEM.Layouts[24] = MEM.Layouts[28]
	MEM.Layouts[23] = MEM.Layouts[31]
	MEM.Layouts[22] = MEM.Layouts[30]
	MEM.Layouts[21] = MEM.Layouts[29]
	MEM.Layouts[20] = MEM.Layouts[28]
	MEM.Layouts[18] = MEM.Layouts[26]
	MEM.Layouts[14] = MEM.Layouts[30]
	MEM.Layouts[9] = MEM.Layouts[25]
	MEM.Layouts[7] = MEM.Layouts[31]
	MEM.Layouts[6] = MEM.Layouts[30]
	MEM.Layouts[4] = MEM.Layouts[28]

	MEM.Layouts[0] = mem.InitConfig(ramSize)
	MEM.Layouts[0].Attach("RAM", 0, RAM, mem.READWRITE)
	MEM.Layouts[0].Accessor("RAM", RAMAccess)
	// MEM.Layouts[0].Attach("FIRST", 0, FIRST, mem.READWRITE)
}

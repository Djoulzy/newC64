package main

import "newC64/mem"

const (
	ramSize     = 65536
	kernalSize  = 8192
	basicSize   = 8192
	ioSize      = 4096
	chargenSize = 4096
)

func main() {
	Mem := mem.Init(5, ramSize)

	RAM := make([]byte, ramSize)
	IO := make([]byte, ioSize)
	KERNAL := mem.LoadROM(kernalSize, "assets/roms/kernal.bin")
	BASIC := mem.LoadROM(basicSize, "assets/roms/basic.bin")
	CHARGEN := mem.LoadROM(chargenSize, "assets/roms/char.bin")

	Mem.Attach("RAM", 0, 0, 0, RAM)
	Mem.Attach("KERNAL", 1, 14, 0, KERNAL)
	Mem.Attach("BASIC", 2, 10, 0, BASIC)
	Mem.Attach("CHARGEN", 3, 13, 0, CHARGEN)
	Mem.Attach("IO", 4, 13, 0, IO)
	Mem.Show()
}

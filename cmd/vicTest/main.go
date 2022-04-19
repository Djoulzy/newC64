package main

import (
	"io/ioutil"
	"newC64/confload"
	"newC64/graphic"
	"newC64/memory"
	"newC64/vic6569"
	"runtime"
)

const (
	ramSize     = 65536
	chargenSize = 4096
	ioSize      = 4096
)

var (
	conf             confload.ConfigData
	mem, io, chargen memory.MEM
	vic              vic6569.VIC
	outputDriver     graphic.Driver
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func LoadData(mem *memory.MEM, file string, memStart uint16) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	for i, val := range content {
		mem.Val[memStart+uint16(i)] = val
	}
	return nil
}

func start() {
	conf.Disassamble = false

	mem.Init(ramSize, "")
	mem.Clear(true)
	io.Init(ioSize, "")
	io.Clear(false)

	chargen.Init(chargenSize, "assets/roms/char.bin")
	LoadData(&mem, "assets/roms/bruce2.bin", 0xE000)

	outputDriver = &graphic.SDLDriver{}
	vic.Init(&mem, &io, &chargen, outputDriver, &conf)
	vic.BankSel = 0
	vic.Write(vic6569.REG_MEM_LOC, 0x78)
	vic.Write(vic6569.REG_CTRL1, 0x3B)
	vic.Write(vic6569.REG_CTRL2, 0x18)
	vic.Write(vic6569.REG_BORDER_COL, 0x0E)
	vic.Write(vic6569.REG_BGCOLOR_0, 0x00)
	vic.Write(vic6569.REG_BGCOLOR_1, 0x01)
	vic.Write(vic6569.REG_BGCOLOR_2, 0x02)
	vic.Write(vic6569.REG_BGCOLOR_3, 0x03)

}

func main() {
	start()
	for {
		vic.Run(false)
		// vic.Stats()
	}
}

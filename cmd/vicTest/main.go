package main

import (
	"io/ioutil"
	"newC64/config"
	"github.com/Djoulzy/emutools/render"
	"github.com/Djoulzy/emutools/mem"
	"newC64/vic6569"
	"runtime"
)

const (
	ramSize     = 65536
	chargenSize = 4096
	ioSize      = 4096
)

var (
	conf             config.ConfigData
	RAM, IO, CHARGEN []byte
	vic              vic6569.VIC
	outputDriver     render.SDL2Driver
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func LoadData(mem []byte, file string, memStart uint16) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	for i, val := range content {
		mem[memStart+uint16(i)] = val
	}
	return nil
}

func start() {
	conf.Disassamble = false

	RAM = make([]byte, ramSize)
	mem.Clear(RAM, 0x100, 0xFF)
	IO = make([]byte, ioSize)
	mem.Clear(IO, 0x100, 0xFF)
	CHARGEN = mem.LoadROM(chargenSize, "assets/roms/char.bin")

	LoadData(RAM, "assets/roms/bruce2.bin", 0xE000)

	outputDriver = render.SDL2Driver{}
	vic.Init(RAM, IO, CHARGEN, &outputDriver, &conf)
	vic.BankSel = 0
	vic.Write(uint16(vic6569.REG_MEM_LOC), 0x78)
	vic.Write(uint16(vic6569.REG_CTRL1), 0x3B)
	vic.Write(uint16(vic6569.REG_CTRL2), 0x18)
	vic.Write(uint16(vic6569.REG_BORDER_COL), 0x0E)
	vic.Write(uint16(vic6569.REG_BGCOLOR_0), 0x00)
	vic.Write(uint16(vic6569.REG_BGCOLOR_1), 0x01)
	vic.Write(uint16(vic6569.REG_BGCOLOR_2), 0x02)
	vic.Write(uint16(vic6569.REG_BGCOLOR_3), 0x03)

}

func main() {
	start()

	go vic.Run(false)
	outputDriver.Run(true)
}

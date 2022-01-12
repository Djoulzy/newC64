package pla906114

import (
	"log"
	"newC64/confload"
	"newC64/memory"
	"os"
	"testing"
)

const (
	ramSize    = 65536
	kernalSize = 8192
	ioSize     = 4096
)

var pla PLA
var conf confload.ConfigData
var mem, io, kernal memory.MEM
var SystemClock uint16

func TestMain(m *testing.M) {
	conf.Disassamble = false
	SystemClock = 0
	var settings byte = 7

	mem.Init(ramSize, "")
	io.Init(ioSize, "")
	kernal.Init(kernalSize, "../assets/roms/kernal.bin")

	pla.Init(&settings)
	pla.Attach(&mem, RAM, 0)
	pla.Attach(&io, IO, IOStart)
	pla.Attach(&kernal, KERNAL, KernalStart)

	os.Exit(m.Run())
}

func TestGetChip(t *testing.T) {
	tables := []struct {
		addr uint16
		res  MemType
	}{
		{0x0100, RAM},
		{0xD000, IO},
		{0xD010, IO},
		{0xE000, KERNAL},
		{0xDC00, IO},
	}

	for _, table := range tables {
		res := pla.getChip(table.addr)
		if res != table.res {
			t.Errorf("GetChip %04X - get: %d - want: %d", table.addr, res, table.res)
		}
	}
	log.Printf("GetChip OK")
}

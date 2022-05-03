package main

import (
	"github.com/Djoulzy/Tools/clog"
)

type ram_accessor struct {
}

func (C *ram_accessor) MRead(mem []byte, translatedAddr uint16) byte {
	// clog.Test("Accessor", "MRead", "Addr: %04X", addr)
	return mem[translatedAddr]
}

func (C *ram_accessor) MWrite(mem []byte, translatedAddr uint16, val byte) {
	switch translatedAddr {
	case 0x0001:
		if LayoutSelector != val&0x1F {
			clog.Test("RAM", "Write", "Layout switch to %08b (%d)", val, val&0x1F)
		}
		LayoutSelector = val & 0x1F
		mem[translatedAddr] = val
	default:
		mem[translatedAddr] = val
	}
}

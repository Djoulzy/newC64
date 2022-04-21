package main

import "newC64/clog"

var IORead_Mapper [](func(uint16) byte)
var IOWrite_Mapper []func(uint16, byte)

func NullRead(addr uint16) byte {
	clog.Trace("NullRead", "Mapper", "addr: %04X - Mapper: %d", addr, IORead_Mapper[addr])
	return MEM.Read(addr)
}

func NullWrite(addr uint16, val byte) {
	clog.Trace("NullWrite", "Mapper", "addr: %04X - Mapper: %d", addr, IORead_Mapper[addr])
	// MEM.Write(addr, val)
}

func fillIOMapper() {
	IORead_Mapper = make([]func(uint16) byte, 4096)
	IOWrite_Mapper = make([]func(uint16, byte), 4096)

	for i := 0; i < 0x0400; i++ {
		IORead_Mapper[i] = vic.Read
		IOWrite_Mapper[i] = vic.Write
	}
	for i := 0x0400; i < 0x0800; i++ { // SID
		IORead_Mapper[i] = NullRead
		IOWrite_Mapper[i] = NullWrite
	}
	for i := 0x0800; i < 0x0C00; i++ { // Color
		IORead_Mapper[i] = MEM.Read
		IOWrite_Mapper[i] = MEM.Write
	}
	for i := 0x0C00; i < 0x0D00; i++ {
		IORead_Mapper[i] = cia1.Read
		IOWrite_Mapper[i] = cia1.Write
	}
	for i := 0x0D00; i < 0x0E00; i++ {
		IORead_Mapper[i] = cia2.Read
		IOWrite_Mapper[i] = cia2.Write
	}
	for i := 0x0E00; i < 0x0F00; i++ { // IO1
		IORead_Mapper[i] = NullRead
		IOWrite_Mapper[i] = NullWrite
	}
	for i := 0x0F00; i < 0x1000; i++ { // IO2
		IORead_Mapper[i] = NullRead
		IOWrite_Mapper[i] = NullWrite
	}
}

package mos6510

import (
	"log"
	"newC64/confload"
	"newC64/memory"
	"newC64/pla906114"
	"os"
	"testing"
)

const (
	ramSize    = 65536
	kernalSize = 8192
	ioSize     = 4096
)

var proc CPU
var pla pla906114.PLA
var conf confload.ConfigData
var mem, io, kernal memory.MEM
var SystemClock uint16

func TestMain(m *testing.M) {
	conf.Disassamble = false
	SystemClock = 0

	mem.Init(ramSize, "")
	io.Init(ioSize, "")
	kernal.Init(kernalSize, "../assets/roms/kernal.bin")

	pla.Init(&mem.Val[1], &conf)
	pla.Attach(&mem, pla906114.RAM, 0)
	pla.Attach(&io, pla906114.IO, pla906114.IOStart)
	pla.Attach(&kernal, pla906114.KERNAL, pla906114.KernalStart)

	proc.Init(&pla, &SystemClock, &conf)
	os.Exit(m.Run())
}

func TestStack(t *testing.T) {
	for i := 0; i <= 0xFF; i++ {
		proc.pushByteStack(byte(i))
	}
	for i := 0xFF; i >= 0; i-- {
		if proc.pullByteStack() != byte(i) {
			t.Errorf("Bad stack operation")
		}
	}

	for i := 0; i <= 0x7F; i++ {
		proc.pushWordStack(uint16(i))
	}
	for i := 0x7F; i >= 0; i-- {
		if proc.pullWordStack() != uint16(i) {
			t.Errorf("Bad stack operation")
		}
	}
	log.Printf("Stack OK")
}

func TestLDA(t *testing.T) {
	proc.inst = mnemonic[0xA9]
	tables := []struct {
		oper byte
		res  byte
		flag byte
	}{
		{0x6E, 0x6E, 0b00000000},
		{0xFF, 0xFF, 0b10000000},
		{0x00, 0x00, 0b00000010},
		{0x81, 0x81, 0b10000000},
	}

	for _, table := range tables {
		proc.S = 0b00000000
		proc.oper = uint16(table.oper)
		proc.lda()
		if proc.S != table.flag {
			t.Errorf("LDA #$%02X - Incorrect status - get: %08b - want: %08b", proc.oper, proc.S, table.flag)
		}
		if proc.A != table.res {
			t.Errorf("LDA #$%02X - Incorrect assignement - get: %02X - want: %02X", proc.oper, proc.A, table.res)
		}
	}
	log.Printf("LDA OK")
}

func TestBNE(t *testing.T) {
	proc.inst = mnemonic[0xD0]
	tables := []struct {
		s    byte
		pc   uint16
		oper byte
		res  uint16
	}{
		{0b00000000, 0xBC16 + uint16(proc.inst.bytes), 0xF9, 0xBC11},
		{0b00000010, 0xBC16 + uint16(proc.inst.bytes), 0xF9, 0xBC18},
	}

	for _, table := range tables {
		proc.PC = table.pc
		proc.S = table.s
		proc.oper = uint16(table.oper)
		proc.bne()
		if proc.PC != table.res {
			t.Errorf("BNE #$%02X - Incorrect status - get: %04X - want: %04X", proc.oper, proc.PC, table.res)
		}
	}
	log.Printf("BNE OK")
}

func TestADC(t *testing.T) {
	proc.inst = mnemonic[0x75] // ZeropageX
	tables := []struct {
		acc     byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x01, 0x04, 0x10, 0b00110000, 0x07, 0b00110000},
		{0x01, 0x04, 0x10, 0b00110001, 0x08, 0b00110000},
		{0xFE, 0x04, 0x10, 0b00110000, 0x04, 0b00110001},
		{0xFE, 0x04, 0x10, 0b00110001, 0x05, 0b00110001},
	}
	proc.ram.Write(0x0014, 0x06)
	for _, table := range tables {
		proc.S = table.flag
		proc.A = table.acc
		proc.X = table.x
		proc.oper = uint16(table.oper)
		proc.adc()
		if proc.A != table.res {
			t.Errorf("A: %02X / ADC $%02X,X - Incorrect result - get: %04X - want: %04X", proc.A, proc.oper, proc.A, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("A: %02X / ADC $%02X,X - Incorrect result Flags - get: %08b - want: %08b", proc.A, proc.oper, proc.S, table.resFlag)
		}
	}

	tables = []struct {
		acc     byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x20, 0x04, 0x10, 0b00110000, 0x2E, 0b00110000},
		{0x01, 0x04, 0x10, 0b00110001, 0x10, 0b00110000},
		{0xA0, 0x04, 0x10, 0b00110000, 0xAE, 0b10110000},
		{0xFE, 0x04, 0x10, 0b00110001, 0x0D, 0b00110001},
	}
	proc.ram.Write(0x0015, 0x02)
	proc.ram.Write(0x0206, 0x0E)
	proc.inst = mnemonic[0x61] // IndirectX
	for _, table := range tables {
		proc.S = table.flag
		proc.A = table.acc
		proc.X = table.x
		proc.oper = uint16(table.oper)
		proc.adc()
		if proc.A != table.res {
			t.Errorf("A: %02X / ADC ($%02X,X) - Incorrect result - get: %04X - want: %04X", proc.A, proc.oper, proc.A, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("A: %02X / ADC ($%02X,X) - Incorrect result Flags - get: %08b - want: %08b", proc.A, proc.oper, proc.S, table.resFlag)
		}
	}

	tables = []struct {
		acc     byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x20, 0x04, 0x14, 0b00110000, 0x2E, 0b00110000},
		{0x01, 0x04, 0x14, 0b00110001, 0x10, 0b00110000},
		{0xA0, 0x04, 0x14, 0b00110000, 0xAE, 0b10110000},
		{0xFE, 0x04, 0x14, 0b00110001, 0x0D, 0b00110001},
	}
	proc.ram.Write(0x0015, 0x02)
	proc.ram.Write(0x020A, 0x0E)
	proc.inst = mnemonic[0x71] // IndirectY
	for _, table := range tables {
		proc.S = table.flag
		proc.A = table.acc
		proc.Y = table.x
		proc.oper = uint16(table.oper)
		proc.adc()
		if proc.A != table.res {
			t.Errorf("A: %02X / ADC ($%02X),Y - Incorrect result - get: %04X - want: %04X", proc.A, proc.oper, proc.A, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("A: %02X / ADC ($%02X),Y - Incorrect result Flags - get: %08b - want: %08b", proc.A, proc.oper, proc.S, table.resFlag)
		}
	}
	log.Printf("ADC OK")
}

func TestSBC(t *testing.T) {
	proc.inst = mnemonic[0xE9] // ZeropageX
	tableIm := []struct {
		acc     byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x03, 0x08, 0b00110000, 0xFA, 0b10110000},
		{0x03, 0x08, 0b00110001, 0xFB, 0b10110000},
	}
	for _, table := range tableIm {
		proc.S = table.flag
		proc.A = table.acc
		proc.oper = uint16(table.oper)
		proc.sbc()
		if proc.A != table.res {
			t.Errorf("A: %02X / SBC #$%02X - Incorrect result - get: %04X - want: %04X", table.acc, proc.oper, proc.A, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("A: %02X / SBC #$%02X - Incorrect result Flags - get: %08b - want: %08b", table.acc, proc.oper, proc.S, table.resFlag)
		}
	}

	proc.inst = mnemonic[0xF5] // ZeropageX
	tables := []struct {
		mem     byte
		acc     byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x06, 0x01, 0x04, 0x10, 0b00110000, 0xFA, 0b10110000},
		{0x06, 0x20, 0x04, 0x10, 0b00110000, 0x19, 0b00110001},
	}

	for _, table := range tables {
		proc.ram.Write(0x0014, table.mem)
		proc.S = table.flag
		proc.A = table.acc
		proc.X = table.x
		proc.oper = uint16(table.oper)
		proc.sbc()
		if proc.A != table.res {
			t.Errorf("A: %02X / SBC $%02X,X - Incorrect result - get: %04X - want: %04X", proc.A, proc.oper, proc.A, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("A: %02X / SBC $%02X,X - Incorrect result Flags - get: %08b - want: %08b", proc.A, proc.oper, proc.S, table.resFlag)
		}
	}

	tables = []struct {
		mem     byte
		acc     byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x06, 0x20, 0x04, 0x10, 0b00110000, 0x11, 0b00110001},
		{0x06, 0x01, 0x04, 0x10, 0b00110001, 0xF3, 0b10110000},
		{0x06, 0xA0, 0x04, 0x10, 0b00110000, 0x91, 0b10110001},
		{0x06, 0xFE, 0x04, 0x10, 0b00110001, 0xF0, 0b10110001},
	}

	proc.inst = mnemonic[0xE1] // IndirectX
	for _, table := range tables {
		proc.ram.Write(0x0014, table.mem)
		proc.ram.Write(0x0015, 0x02)
		proc.ram.Write(0x0206, 0x0E)
		proc.S = table.flag
		proc.A = table.acc
		proc.X = table.x
		proc.oper = uint16(table.oper)
		proc.sbc()
		if proc.A != table.res {
			t.Errorf("A: %02X / SBC ($%02X,X) - Incorrect result - get: %04X - want: %04X", proc.A, proc.oper, proc.A, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("A: %02X / SBC ($%02X,X) - Incorrect result Flags - get: %08b - want: %08b", proc.A, proc.oper, proc.S, table.resFlag)
		}
	}

	tables = []struct {
		mem     byte
		acc     byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x06, 0x20, 0x04, 0x14, 0b00110000, 0x11, 0b00110001},
		{0x06, 0x01, 0x04, 0x14, 0b00110001, 0xF3, 0b10110000},
		{0x06, 0xA0, 0x04, 0x14, 0b00110000, 0x91, 0b10110001},
		{0x06, 0xFE, 0x04, 0x14, 0b00110001, 0xF0, 0b10110001},
	}
	proc.ram.Write(0x0015, 0x02)
	proc.ram.Write(0x020A, 0x0E)
	proc.inst = mnemonic[0xF1] // IndirectY
	for _, table := range tables {
		proc.ram.Write(0x0014, table.mem)
		proc.ram.Write(0x0015, 0x02)
		proc.ram.Write(0x0206, 0x0E)
		proc.S = table.flag
		proc.A = table.acc
		proc.Y = table.x
		proc.oper = uint16(table.oper)
		proc.sbc()
		if proc.A != table.res {
			t.Errorf("A: %02X / SBC ($%02X),Y - Incorrect result - get: %04X - want: %04X", proc.A, proc.oper, proc.A, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("A: %02X / SBC ($%02X),Y - Incorrect result Flags - get: %08b - want: %08b", proc.A, proc.oper, proc.S, table.resFlag)
		}
	}
	log.Printf("SBC OK")
}

func TestCMP(t *testing.T) {
	tables := []struct {
		acc  byte
		oper byte
		flag byte
	}{
		{0x50, 0x20, 0b00110001},
		{0xF0, 0x20, 0b10110001},
		{0x00, 0x20, 0b10110000},
		{0x20, 0x20, 0b00110011},
		{0x01, 0x20, 0b10110000},
		{0x00, 0x00, 0b00110011},
		{0xFF, 0xFF, 0b00110011},
	}

	proc.inst = mnemonic[0xC9]
	for _, table := range tables {
		proc.S = 0b00110000
		proc.A = table.acc
		proc.oper = uint16(table.oper)
		proc.cmp()
		if proc.S != table.flag {
			t.Errorf("LDA #$%02X;CMP #$%02X - Incorrect status - get: %08b - want: %08b", proc.A, proc.oper, proc.S, table.flag)
		}
	}

	tables = []struct {
		acc  byte
		oper byte
		flag byte
	}{
		{0x50, 0xC1, 0b00110000},
		{0xF0, 0xC1, 0b00110001},
		{0x00, 0xC1, 0b00110000},
		{0x20, 0xC1, 0b00110000},
		{0xEE, 0xC1, 0b00110011},
		{0xFF, 0xC1, 0b00110001},
	}

	proc.inst = mnemonic[0xD1]
	proc.ram.Write(0x0408, 0xEE)
	proc.ram.Write(0xC1, 0x00)
	proc.ram.Write(0xC2, 0x04)
	for _, table := range tables {
		proc.S = 0b00110000
		proc.Y = 0x08
		proc.A = table.acc
		proc.oper = uint16(table.oper)
		proc.cmp()
		if proc.S != table.flag {
			t.Errorf("LDA #$%02X;CMP ($%02X),Y - Incorrect status - get: %08b - want: %08b", proc.A, proc.oper, proc.S, table.flag)
		}
	}
	log.Printf("CMP OK")
}

func TestROR(t *testing.T) {
	proc.inst = mnemonic[0x76] // ZeropageX
	tables := []struct {
		val     byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x06, 0x04, 0x10, 0b00110000, 0x03, 0b00110000},
		{0x06, 0x04, 0x10, 0b00110001, 0x83, 0b10110000},
	}

	for _, table := range tables {
		proc.ram.Write(0x0014, table.val)
		proc.S = table.flag
		proc.X = table.x
		proc.oper = uint16(table.oper)
		proc.ror()
		res := proc.ram.Read(0x0014)
		if res != table.res {
			t.Errorf("Val: $%02X / ROR $%02X,X - Incorrect result - get: %02X - want: %02X", table.val, proc.oper, res, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("Val: $%02X / ROR $%02X,X - Incorrect result Flags - get: %08b - want: %08b", table.val, proc.oper, proc.S, table.resFlag)
		}
	}
	log.Printf("ROR OK")
}

func TestROL(t *testing.T) {
	proc.inst = mnemonic[0x76] // ZeropageX
	tables := []struct {
		val     byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x06, 0x04, 0x10, 0b00110000, 0x0C, 0b00110000},
		{0x06, 0x04, 0x10, 0b00110001, 0x0D, 0b00110000},
		{0x80, 0x04, 0x10, 0b00110001, 0x01, 0b00110001},
		{0xF0, 0x04, 0x10, 0b00110001, 0xE1, 0b10110001},
		{0xF0, 0x04, 0x10, 0b00110000, 0xE0, 0b10110001},
	}

	for _, table := range tables {
		proc.ram.Write(0x0014, table.val)
		proc.S = table.flag
		proc.X = table.x
		proc.oper = uint16(table.oper)
		proc.rol()
		res := proc.ram.Read(0x0014)
		if res != table.res {
			t.Errorf("Val: $%02X / ROL $%02X,X - Incorrect result - get: %02X - want: %02X", table.val, proc.oper, res, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("Val: $%02X / ROL $%02X,X - Incorrect result Flags - get: %08b - want: %08b", table.val, proc.oper, proc.S, table.resFlag)
		}
	}
	log.Printf("ROL OK")
}

func TestLSR(t *testing.T) {
	proc.inst = mnemonic[0x56] // ZeropageX
	tables := []struct {
		val     byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x80, 0x04, 0x10, 0b00110000, 0x40, 0b00110000},
		{0x0F, 0x04, 0x10, 0b00110000, 0x07, 0b00110001},
		{0x0F, 0x04, 0x10, 0b00110001, 0x07, 0b00110001},
		{0x80, 0x04, 0x10, 0b00110001, 0x40, 0b00110000},
		{0xFF, 0x04, 0x10, 0b00110001, 0x7F, 0b00110001},
	}

	for _, table := range tables {
		proc.ram.Write(0x0014, table.val)
		proc.S = table.flag
		proc.X = table.x
		proc.oper = uint16(table.oper)
		proc.lsr()
		res := proc.ram.Read(0x0014)
		if res != table.res {
			t.Errorf("Val: $%02X / LSR $%02X,X - Incorrect result - get: %02X - want: %02X", table.val, proc.oper, res, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("Val: $%02X / LSR $%02X,X - Incorrect result Flags - get: %08b - want: %08b", table.val, proc.oper, proc.S, table.resFlag)
		}
	}
	log.Printf("LSR OK")
}

func TestASL(t *testing.T) {
	proc.inst = mnemonic[0x16] // ZeropageX
	tables := []struct {
		val     byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x80, 0x04, 0x10, 0b00110000, 0x00, 0b00110011},
		{0x7F, 0x04, 0x10, 0b00110000, 0xFE, 0b10110000},
		{0x7F, 0x04, 0x10, 0b00110001, 0xFE, 0b10110000},
		{0x80, 0x04, 0x10, 0b00110001, 0x00, 0b00110011},
		{0xFF, 0x04, 0x10, 0b00110001, 0xFE, 0b10110001},
	}

	for _, table := range tables {
		proc.ram.Write(0x0014, table.val)
		proc.S = table.flag
		proc.X = table.x
		proc.oper = uint16(table.oper)
		proc.asl()
		res := proc.ram.Read(0x0014)
		if res != table.res {
			t.Errorf("Val: $%02X / ASL $%02X,X - Incorrect result - get: %02X - want: %02X", table.val, proc.oper, res, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("Val: $%02X / ASL $%02X,X - Incorrect result Flags - get: %08b - want: %08b", table.val, proc.oper, proc.S, table.resFlag)
		}
	}
	log.Printf("ASL OK")
}

func TestEOR(t *testing.T) {
	// LDA #$80
	// STA $14
	// LDX #$04
	// CLC
	// LDA #$11
	// EOR $10,X
	proc.inst = mnemonic[0x55] // ZeropageX
	tables := []struct {
		val     byte
		a       byte
		x       byte
		oper    byte
		flag    byte
		res     byte
		resFlag byte
	}{
		{0x80, 0x11, 0x04, 0x10, 0b00110000, 0x91, 0b10110000},
		{0x80, 0x80, 0x04, 0x10, 0b00110000, 0x00, 0b00110010},
		{0x80, 0x0F, 0x04, 0x10, 0b00110001, 0x8F, 0b10110001},
		{0x80, 0xFF, 0x04, 0x10, 0b00110001, 0x7F, 0b00110001},
		{0x80, 0x00, 0x04, 0x10, 0b00110001, 0x80, 0b10110001},
	}

	for _, table := range tables {
		proc.ram.Write(0x0014, table.val)
		proc.S = table.flag
		proc.X = table.x
		proc.A = table.a
		proc.oper = uint16(table.oper)
		proc.eor()
		if proc.A != table.res {
			t.Errorf("LDA #$%02X / EOR $%02X,X - Incorrect result - get: %02X - want: %02X", table.a, proc.oper, proc.A, table.res)
		}
		if proc.S != table.resFlag {
			t.Errorf("LDA #$%02X / EOR $%02X,X - Incorrect result Flags - get: %08b - want: %08b", table.a, proc.oper, proc.S, table.resFlag)
		}
	}
	log.Printf("EOR OK")
}

func TestBIT(t *testing.T) {
	// LDA #$80
	// STA $14
	// CLC
	// LDA #$11
	// BIT $14
	proc.inst = mnemonic[0x24] // Zeropage
	tables := []struct {
		val     byte
		a       byte
		flag    byte
		resFlag byte
	}{
		{0x80, 0x11, 0b00110000, 0b10110010},
		{0x80, 0x80, 0b00110000, 0b10110000},
		{0x80, 0x0F, 0b00110001, 0b10110011},
		{0x80, 0xFF, 0b00110001, 0b10110001},
		{0x80, 0x00, 0b00110011, 0b10110011},
	}

	for _, table := range tables {
		proc.ram.Write(0x0014, table.val)
		proc.S = table.flag
		proc.A = table.a
		proc.oper = 0x14
		proc.bit()
		if proc.S != table.resFlag {
			t.Errorf("LDA #$%02X / LDA #$%02X -> BIT $%02X - Incorrect result Flags - get: %08b - want: %08b", table.val, table.a, proc.oper, proc.S, table.resFlag)
		}
	}
	log.Printf("BIT OK")
}

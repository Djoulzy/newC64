package mos6510

import (
	"newC64/confload"
	"newC64/pla906114"
	"os"
	"testing"
)

const (
	memorySize = 65536
	kernalSize = 8192
)

var proc CPU
var pla pla906114.PLA
var conf confload.ConfigData

func TestMain(m *testing.M) {
	conf.Disassamble = false

	mem := make([]byte, memorySize)
	kernal := make([]byte, kernalSize)

	pla.Init()
	pla.Attach(mem, pla906114.RAM, 0)
	pla.Attach(kernal, pla906114.KERNAL, pla906114.KernalStart)

	proc.Init(&pla, &conf)
	os.Exit(m.Run())
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
}

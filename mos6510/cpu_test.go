package mos6510

import (
	"log"
	"newC64/confload"
	"newC64/mem"
	"os"
	"testing"
)

const (
	ramSize    = 65536
	kernalSize = 8192
	ioSize     = 4096
)

type TestData struct {
	proc    *CPU
	inst    byte
	mem     byte
	memVal  byte
	memDest byte
	acc     byte
	x       byte
	y       byte
	pc      uint16
	oper    byte
	flag    byte
	res     uint16
	resFlag byte
}

type TestSuite struct {
	proc *CPU
	inst byte
	data []TestData
}

func getAddrName(a addressing) string {
	switch a {
	case implied:
		return "implied"
	case immediate:
		return "immediate"
	case relative:
		return "relative"
	case zeropage:
		return "zeropage"
	case zeropageX:
		return "zeropageX"
	case zeropageY:
		return "zeropageY"
	case absolute:
		return "absolute"
	case absoluteX:
		return "absoluteX"
	case absoluteY:
		return "absoluteY"
	case indirect:
		return "indirect"
	case indirectX:
		return "indirectX"
	case indirectY:
		return "indirectY"
	case Branching:
		return "Branching"
	case CrossPage:
		return "CrossPage"
	}
	return "Unknown"
}

func (TS *TestSuite) Add(td TestData) {
	td.proc = TS.proc
	td.inst = TS.inst
	TS.data = append(TS.data, td)
}

func (TD *TestData) run() {
	proc.Inst = mnemonic[TD.inst]
	proc.PC = TD.pc + uint16(proc.Inst.bytes)
	proc.S = TD.flag
	proc.A = TD.acc
	proc.X = TD.x
	proc.Y = TD.y
	proc.oper = uint16(TD.oper)
	for {
		cycle := proc.Inst.Cycles
		proc.Inst.action()
		if cycle == proc.Inst.Cycles {
			break
		}
	}
}

func (TD *TestData) checkBit(t *testing.T, val1, val2 byte, name string) bool {
	if val1 != val2 {
		t.Errorf("%s %s - Incorrect %s - get: %08b - want: %08b", mnemonic[TD.inst].name, getAddrName(mnemonic[TD.inst].addr), name, val1, val2)
		return false
	}
	return true
}

func (TD *TestData) checkByte(t *testing.T, val1, val2 byte, name string) bool {
	if val1 != val2 {
		t.Errorf("%s %s - Incorrect %s - get: %02X - want: %02X", mnemonic[TD.inst].name, getAddrName(mnemonic[TD.inst].addr), name, val1, val2)
		return false
	}
	return true
}

func (TD *TestData) checkWord(t *testing.T, val1, val2 uint16, name string) bool {
	if val1 != val2 {
		t.Errorf("%s %s - Incorrect %s - get: %04X - want: %04X", mnemonic[TD.inst].name, getAddrName(mnemonic[TD.inst].addr), name, val1, val2)
		return false
	}
	return true
}

func finalize(name string, allGood bool) {
	if allGood {
		log.Printf("%s OK", name)
	} else {
		log.Printf("%s %c[1;31mECHEC%c[0m", name, 27, 27)
	}
}

var proc CPU
var BankSel byte
var MEM mem.BANK
var conf confload.ConfigData
var RAM, IO, KERNAL []byte
var SystemClock uint16

func TestMain(m *testing.M) {
	conf.Disassamble = false
	SystemClock = 0

	RAM = make([]byte, ramSize)
	IO = make([]byte, ioSize)
	KERNAL = mem.LoadROM(kernalSize, "../assets/roms/kernal.bin")

	BankSel = 0
	MEM = mem.InitBanks(1, &BankSel)

	MEM.Layouts[0] = mem.InitConfig(3, ramSize)
	MEM.Layouts[0].Attach("RAM", 0, 0, RAM, mem.READWRITE)
	MEM.Layouts[0].Attach("IO", 1, 13, IO, mem.READWRITE)
	MEM.Layouts[0].Attach("KERNAL", 2, 14, KERNAL, mem.READONLY)
	MEM.Layouts[0].Show()

	proc.Init(&MEM, &SystemClock, &conf)
	os.Exit(m.Run())
}

func TestStack(t *testing.T) {
	var allGood bool = true
	mem.Clear(RAM)
	for i := 0; i <= 0xFF; i++ {
		proc.pushByteStack(byte(i))
	}
	for i := 0xFF; i >= 0; i-- {
		if proc.pullByteStack() != byte(i) {
			t.Errorf("Bad stack operation")
			allGood = false
		}
	}

	for i := 0; i <= 0x7F; i++ {
		proc.pushWordStack(uint16(i))
	}
	for i := 0x7F; i >= 0; i-- {
		if proc.pullWordStack() != uint16(i) {
			t.Errorf("Bad stack operation")
			allGood = false
		}
	}
	finalize("Stack", allGood)
}

func TestLDA(t *testing.T) {
	var allGood bool = true
	mem.Clear(RAM)

	ts := TestSuite{proc: &proc, inst: 0xA9}
	ts.Add(TestData{oper: 0x6E, res: 0x6E, flag: 0b00100000, resFlag: 0b00100000})
	ts.Add(TestData{oper: 0xFF, res: 0xFF, flag: 0b00100000, resFlag: 0b10100000})
	ts.Add(TestData{oper: 0x00, res: 0x00, flag: 0b00100000, resFlag: 0b00100010})
	ts.Add(TestData{oper: 0x81, res: 0x81, flag: 0b00100000, resFlag: 0b10100000})

	for _, table := range ts.data {
		table.run()
		allGood = allGood && table.checkBit(t, proc.S, table.resFlag, "Status Flag")
		allGood = allGood && table.checkByte(t, proc.A, byte(table.res), "Assignement")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

func TestBNE(t *testing.T) {
	var allGood bool = true
	mem.Clear(RAM)
	ts := TestSuite{proc: &proc, inst: 0xD0}

	ts.Add(TestData{flag: 0b00000000, pc: 0xBC16, oper: 0xF9, res: 0xBC11})
	ts.Add(TestData{flag: 0b00000010, pc: 0xBC16, oper: 0xF9, res: 0xBC18})

	for _, table := range ts.data {
		table.run()
		allGood = allGood && table.checkWord(t, proc.PC, table.res, "Address")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

func TestADC(t *testing.T) {
	// LDA #$06
	// STA $0014
	// LDA #$02
	// STA $0015

	// LDA #$0E
	// STA $020A

	// LDY #$04
	// LDA #$20
	// CLC
	// ADC ($14),Y
	var allGood bool = true
	mem.Clear(RAM)

	ts := TestSuite{proc: &proc, inst: 0x75}
	ts.Add(TestData{acc: 0x01, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x07, resFlag: 0b00110000})
	ts.Add(TestData{acc: 0x01, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x08, resFlag: 0b00110000})
	ts.Add(TestData{acc: 0xFE, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x04, resFlag: 0b00110001})
	ts.Add(TestData{acc: 0xFE, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x05, resFlag: 0b00110001})

	ts.inst = 0x69 // immediate
	ts.Add(TestData{acc: 0x78, x: 0x04, oper: 0x80, flag: 0b00110000, res: 0xF8, resFlag: 0b10110000})
	ts.Add(TestData{acc: 0x80, x: 0x04, oper: 0x12, flag: 0b00111000, res: 0x92, resFlag: 0b10111000})
	ts.Add(TestData{acc: 0x58, x: 0x04, oper: 0x46, flag: 0b00111001, res: 0x05, resFlag: 0b01111001})
	ts.Add(TestData{acc: 0x99, x: 0x04, oper: 0x01, flag: 0b00111000, res: 0x00, resFlag: 0b00111011})

	ts.inst = 0x61
	ts.Add(TestData{acc: 0x20, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x2E, resFlag: 0b00110000})
	ts.Add(TestData{acc: 0x01, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x10, resFlag: 0b00110000})
	ts.Add(TestData{acc: 0xA0, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0xAE, resFlag: 0b10110000})
	ts.Add(TestData{acc: 0xFE, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x0D, resFlag: 0b00110001})

	ts.inst = 0x71
	ts.Add(TestData{acc: 0x20, x: 0x04, oper: 0x14, flag: 0b00110000, res: 0x2E, resFlag: 0b00110000})
	ts.Add(TestData{acc: 0x01, x: 0x04, oper: 0x14, flag: 0b00110001, res: 0x10, resFlag: 0b00110000})
	ts.Add(TestData{acc: 0xA0, x: 0x04, oper: 0x14, flag: 0b00110000, res: 0xAE, resFlag: 0b10110000})
	ts.Add(TestData{acc: 0xFE, x: 0x04, oper: 0x14, flag: 0b00110001, res: 0x0D, resFlag: 0b00110001})

	proc.ram.Write(0x0014, 0x06)
	proc.ram.Write(0x0015, 0x02)
	proc.ram.Write(0x0206, 0x0E)
	proc.ram.Write(0x020A, 0x0E)
	for _, table := range ts.data {
		table.run()
		allGood = allGood && table.checkByte(t, proc.A, byte(table.res), "Result")
		allGood = allGood && table.checkBit(t, proc.S, table.resFlag, "Flags")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

func TestSBC(t *testing.T) {
	// LDA #$06
	// STA $0014
	// LDA #$02
	// STA $0015

	// LDA #$0E
	// STA $020A

	// LDY #$04
	// LDA #$fe
	// SEC
	// SBC ($14),Y
	var allGood bool = true
	mem.Clear(RAM)
	ts := TestSuite{proc: &proc, inst: 0xE9} // Immediate
	ts.Add(TestData{acc: 0x03, oper: 0x08, flag: 0b00110000, res: 0xFA, resFlag: 0b10110000})
	ts.Add(TestData{acc: 0x03, oper: 0x08, flag: 0b00110001, res: 0xFB, resFlag: 0b10110000})
	ts.Add(TestData{acc: 0x58, oper: 0x46, flag: 0b00111000, res: 0x11, resFlag: 0b00111001})

	ts.inst = 0xF5
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0x01, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0xFA, resFlag: 0b10110000})
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0x20, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x19, resFlag: 0b00110001})

	ts.inst = 0xE1
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0x20, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x11, resFlag: 0b00110001})
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0x01, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0xF3, resFlag: 0b10110000})
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0xA0, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x91, resFlag: 0b10110001})
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0xFE, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0xF0, resFlag: 0b10110001})

	ts.inst = 0xF1
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0x20, y: 0x04, oper: 0x14, flag: 0b00110000, res: 0x11, resFlag: 0b00110001})
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0x01, y: 0x04, oper: 0x14, flag: 0b00110001, res: 0xF3, resFlag: 0b10110000})
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0xA0, y: 0x04, oper: 0x14, flag: 0b00110000, res: 0x91, resFlag: 0b10110001})
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0xFE, y: 0x04, oper: 0x14, flag: 0b00110001, res: 0xF0, resFlag: 0b10110001})
	ts.Add(TestData{mem: 0x06, memVal: 0x08, acc: 0x03, y: 0x04, oper: 0x14, flag: 0b00110001, res: 0xFB, resFlag: 0b10110000})
	ts.Add(TestData{mem: 0x06, memVal: 0x08, acc: 0x03, y: 0x04, oper: 0x14, flag: 0b00110000, res: 0xFA, resFlag: 0b10110000})
	ts.Add(TestData{mem: 0x06, memVal: 0x0E, acc: 0xFE, y: 0x04, oper: 0x14, flag: 0b00110011, res: 0xF0, resFlag: 0b10110001})

	for _, table := range ts.data {
		proc.ram.Write(0x0014, table.mem)
		proc.ram.Write(0x0015, 0x02)
		proc.ram.Write(0x0206, table.memVal)
		proc.ram.Write(0x020A, table.memVal)
		table.run()
		allGood = allGood && table.checkByte(t, proc.A, byte(table.res), "Result")
		allGood = allGood && table.checkBit(t, proc.S, table.resFlag, "Flags")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

func TestCMP(t *testing.T) {
	var allGood bool = true
	mem.Clear(RAM)

	ts := TestSuite{proc: &proc, inst: 0xC9}
	ts.Add(TestData{acc: 0x50, oper: 0x20, flag: 0b00110000, resFlag: 0b00110001})
	ts.Add(TestData{acc: 0xF0, oper: 0x20, flag: 0b00110000, resFlag: 0b10110001})
	ts.Add(TestData{acc: 0x00, oper: 0x20, flag: 0b00110000, resFlag: 0b10110000})
	ts.Add(TestData{acc: 0x20, oper: 0x20, flag: 0b00110000, resFlag: 0b00110011})
	ts.Add(TestData{acc: 0x01, oper: 0x20, flag: 0b00110000, resFlag: 0b10110000})
	ts.Add(TestData{acc: 0x00, oper: 0x00, flag: 0b00110000, resFlag: 0b00110011})
	ts.Add(TestData{acc: 0xFF, oper: 0xFF, flag: 0b00110000, resFlag: 0b00110011})

	ts.inst = 0xD1
	ts.Add(TestData{acc: 0x50, y: 0x08, oper: 0xC1, flag: 0b00110000, resFlag: 0b00110000})
	ts.Add(TestData{acc: 0xF0, y: 0x08, oper: 0xC1, flag: 0b00110000, resFlag: 0b00110001})
	ts.Add(TestData{acc: 0x00, y: 0x08, oper: 0xC1, flag: 0b00110000, resFlag: 0b00110000})
	ts.Add(TestData{acc: 0x20, y: 0x08, oper: 0xC1, flag: 0b00110000, resFlag: 0b00110000})
	ts.Add(TestData{acc: 0xEE, y: 0x08, oper: 0xC1, flag: 0b00110000, resFlag: 0b00110011})
	ts.Add(TestData{acc: 0xFF, y: 0x08, oper: 0xC1, flag: 0b00110000, resFlag: 0b00110001})

	for _, table := range ts.data {
		proc.ram.Write(0x0408, 0xEE)
		proc.ram.Write(0xC1, 0x00)
		proc.ram.Write(0xC2, 0x04)
		table.run()
		allGood = allGood && table.checkBit(t, proc.S, table.resFlag, "Flags")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

func TestROR(t *testing.T) {
	var allGood bool = true
	mem.Clear(RAM)
	ts := TestSuite{proc: &proc, inst: 0x76}

	ts.Add(TestData{mem: 0x06, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x03, resFlag: 0b00110000})
	ts.Add(TestData{mem: 0x06, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x83, resFlag: 0b10110000})

	for _, table := range ts.data {
		proc.ram.Write(0x0014, table.mem)
		table.run()
		allGood = allGood && table.checkByte(t, proc.ram.Read(0x0014), byte(table.res), "Result")
		allGood = allGood && table.checkBit(t, proc.S, table.resFlag, "Flags")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

func TestROL(t *testing.T) {
	var allGood bool = true
	mem.Clear(RAM)
	ts := TestSuite{proc: &proc, inst: 0x36}
	ts.Add(TestData{mem: 0x06, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x0C, resFlag: 0b00110000})
	ts.Add(TestData{mem: 0x06, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x0D, resFlag: 0b00110000})
	ts.Add(TestData{mem: 0x80, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x01, resFlag: 0b00110001})
	ts.Add(TestData{mem: 0xF0, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0xE1, resFlag: 0b10110001})
	ts.Add(TestData{mem: 0xF0, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0xE0, resFlag: 0b10110001})

	for _, table := range ts.data {
		proc.ram.Write(0x0014, table.mem)
		table.run()
		allGood = allGood && table.checkByte(t, proc.ram.Read(0x0014), byte(table.res), "Result")
		allGood = allGood && table.checkBit(t, proc.S, table.resFlag, "Flags")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

func TestLSR(t *testing.T) {
	var allGood bool = true
	mem.Clear(RAM)
	ts := TestSuite{proc: &proc, inst: 0x56}
	ts.Add(TestData{mem: 0x80, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x40, resFlag: 0b00110000})
	ts.Add(TestData{mem: 0x0F, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x07, resFlag: 0b00110001})
	ts.Add(TestData{mem: 0x0F, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x07, resFlag: 0b00110001})
	ts.Add(TestData{mem: 0x80, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x40, resFlag: 0b00110000})
	ts.Add(TestData{mem: 0xFF, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x7F, resFlag: 0b00110001})

	for _, table := range ts.data {
		proc.ram.Write(0x0014, table.mem)
		table.run()
		allGood = allGood && table.checkByte(t, proc.ram.Read(0x0014), byte(table.res), "Result")
		allGood = allGood && table.checkBit(t, proc.S, table.resFlag, "Flags")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

func TestASL(t *testing.T) {
	var allGood bool = true
	mem.Clear(RAM)
	ts := TestSuite{proc: &proc, inst: 0x0E}
	ts.Add(TestData{mem: 0xF1, x: 0x04, memDest: 0x19, oper: 0x0019, flag: 0b00100101, res: 0xE2, resFlag: 0b10100101})

	ts.inst = 0x16
	ts.Add(TestData{mem: 0x80, x: 0x04, memDest: 0x14, oper: 0x10, flag: 0b00110000, res: 0x00, resFlag: 0b00110011})
	ts.Add(TestData{mem: 0x7F, x: 0x04, memDest: 0x14, oper: 0x10, flag: 0b00110000, res: 0xFE, resFlag: 0b10110000})
	ts.Add(TestData{mem: 0x7F, x: 0x04, memDest: 0x14, oper: 0x10, flag: 0b00110001, res: 0xFE, resFlag: 0b10110000})
	ts.Add(TestData{mem: 0x80, x: 0x04, memDest: 0x14, oper: 0x10, flag: 0b00110001, res: 0x00, resFlag: 0b00110011})
	ts.Add(TestData{mem: 0xFF, x: 0x04, memDest: 0x14, oper: 0x10, flag: 0b00110001, res: 0xFE, resFlag: 0b10110001})

	for _, table := range ts.data {
		proc.ram.Write(uint16(table.memDest), table.mem)
		table.run()
		allGood = allGood && table.checkByte(t, proc.ram.Read(uint16(table.memDest)), byte(table.res), "Result")
		allGood = allGood && table.checkBit(t, proc.S, table.resFlag, "Flags")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

func TestEOR(t *testing.T) {
	var allGood bool = true
	mem.Clear(RAM)
	ts := TestSuite{proc: &proc, inst: 0x55}
	ts.Add(TestData{mem: 0x80, acc: 0x11, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x91, resFlag: 0b10110000})
	ts.Add(TestData{mem: 0x80, acc: 0x80, x: 0x04, oper: 0x10, flag: 0b00110000, res: 0x00, resFlag: 0b00110010})
	ts.Add(TestData{mem: 0x80, acc: 0x0F, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x8F, resFlag: 0b10110001})
	ts.Add(TestData{mem: 0x80, acc: 0xFF, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x7F, resFlag: 0b00110001})
	ts.Add(TestData{mem: 0x80, acc: 0x00, x: 0x04, oper: 0x10, flag: 0b00110001, res: 0x80, resFlag: 0b10110001})

	for _, table := range ts.data {
		proc.ram.Write(0x0014, table.mem)
		table.run()
		allGood = allGood && table.checkByte(t, proc.A, byte(table.res), "Result")
		allGood = allGood && table.checkBit(t, proc.S, table.resFlag, "Flags")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

func TestBIT(t *testing.T) {
	var allGood bool = true
	// LDA #$80
	// STA $14
	// CLC
	// LDA #$11
	// BIT $14
	mem.Clear(RAM)
	ts := TestSuite{proc: &proc, inst: 0x24}

	ts.Add(TestData{mem: 0x80, acc: 0x11, oper: 0x14, flag: 0b00110000, resFlag: 0b10110010})
	ts.Add(TestData{mem: 0x80, acc: 0x80, oper: 0x14, flag: 0b00110000, resFlag: 0b10110000})
	ts.Add(TestData{mem: 0x80, acc: 0x0F, oper: 0x14, flag: 0b00110001, resFlag: 0b10110011})
	ts.Add(TestData{mem: 0x80, acc: 0xFF, oper: 0x14, flag: 0b00110001, resFlag: 0b10110001})
	ts.Add(TestData{mem: 0x80, acc: 0x00, oper: 0x14, flag: 0b00110011, resFlag: 0b10110011})

	for _, table := range ts.data {
		proc.ram.Write(0x0014, table.mem)
		table.run()
		allGood = allGood && table.checkBit(t, proc.S, table.resFlag, "Flags")
	}
	finalize(mnemonic[ts.inst].name, allGood)
}

package main

import (
	"fmt"
	"log"
	"newC64/cia6526"
	"newC64/config"
	"newC64/vic6569"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Djoulzy/emutools/render"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/Tools/confload"
	"github.com/Djoulzy/emutools/mem"
	"github.com/Djoulzy/emutools/mos6510"

	"github.com/mattn/go-tty"
)

const (
	ramSize     = 65536
	kernalSize  = 8192
	basicSize   = 8192
	ioSize      = 4096
	chargenSize = 4096
	cartSize    = 8192

	nbMemLayout = 32

	Stopped = 0
	Paused  = 1
	Running = 2
)

var (
	conf = &config.ConfigData{}

	cpu  mos6510.CPU
	cia1 cia6526.CIA
	cia2 cia6526.CIA

	RAM      []byte
	IO       []byte
	KERNAL   []byte
	BASIC    []byte
	CHARGEN  []byte
	CART_LO  []byte
	CART_HI  []byte
	MEM      mem.BANK
	IOAccess mem.MEMAccess

	vic vic6569.VIC

	outputDriver render.SDL2Driver
	cpuTurn      bool
	run          bool
	execInst     sync.Mutex
)

// func init() {
// 	// This is needed to arrange that main() runs on main thread.
// 	// See documentation for functions that are only allowed to be called from the main thread.
// 	runtime.LockOSThread()
// }

func setup() {
	// ROMs & RAM Setup
	RAM = make([]byte, ramSize)
	IO = make([]byte, ioSize)
	CART_LO = make([]byte, cartSize)
	CART_HI = make([]byte, cartSize)
	KERNAL = mem.LoadROM(kernalSize, "assets/roms/kernal.bin")
	BASIC = mem.LoadROM(basicSize, "assets/roms/basic.bin")
	CHARGEN = mem.LoadROM(chargenSize, "assets/roms/char.bin")

	mem.Clear(RAM)
	mem.Clear(IO)

	RAM[0x0001] = 0x1F
	MEM = mem.InitBanks(nbMemLayout, &RAM[0x0001])
	// var test byte = 31
	// MEM = mem.InitBanks(nbMemLayout, &test)
	IOAccess = &accessor{}
	fillIOMapper()

	// MEM Setup
	memLayouts()

	outputDriver = render.SDL2Driver{}
	vic.Init(RAM, IO, CHARGEN, &outputDriver, conf)

	// CPU Setup
	cpu.Init(&MEM)

	cia1.Init("CIA1", IO[0x0C00:0x0C00+0x0200], &vic.SystemClock)
	outputDriver.SetKeyboardLine(&cia1.InputLine)
	cia2.Init("CIA2", IO[0x0D00:0x0D00+0x0200], &vic.SystemClock)

	vic.IRQ_Pin = &cpu.IRQ_pin
	cia1.Signal_Pin = &cpu.IRQ_pin
	cia2.Signal_Pin = &cpu.NMI_pin
	cia2.VICBankSelect = &vic.BankSel
}

func input() {
	dumpAddr := ""
	var keyb *tty.TTY
	keyb, _ = tty.Open()

	for {
		r, _ := keyb.ReadRune()
		switch r {
		case 'v':
			cia1.Stats()
			vic.Stats()
		case 's':
			Disassamble()
			MEM.DumpStack(cpu.SP)
		case 'z':
			Disassamble()
			MEM.Dump(0)
		case 'x':
			// DumpMem(&pla, "memDump.bin")
		case 'r':
			conf.Disassamble = false
			run = true
			execInst.Unlock()
		case 'l':
			// LoadPRG(&pla, "./prg/GARDEN.prg")
			LoadPRG(&MEM, conf.LoadPRG)
			// addr, _ := LoadPRG(mem.Val, conf.LoadPRG)
			// cpu.GoTo(addr)
		case ' ':
			if run {
				conf.Disassamble = true
				run = false
			} else {
				execInst.Unlock()
			}
			// fmt.Printf("\n(s) Stack Dump - (z) Zero Page - (r) Run - (sp) Pause / unpause > ")
		case 'w':
			fmt.Printf("\nFill Color RAM")
			for i := 0xD800; i < 0xDC00; i++ {
				MEM.Write(uint16(i), 0)
			}
			// for i := 0x0800; i < 0x0C00; i++ {
			// 	IO[uint16(i)] = 0
			// }
		case 'q':
			cpu.DumpStats()
			os.Exit(0)
		default:
			dumpAddr += string(r)
			fmt.Printf("%c", r)
			if len(dumpAddr) == 4 {
				hx, _ := strconv.ParseInt(dumpAddr, 16, 64)
				MEM.Dump(uint16(hx))
				if hx < 0x4000 {
					vic.Dump(uint16(hx))
				}
				dumpAddr = ""
			}
		}

	}
}

func Disassamble() {
	// fmt.Printf("\n%s %s", vic.Disassemble(), cpu.Disassemble())
	fmt.Printf("%s\n", cpu.Trace())
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Now().Sub(start)
	log.Printf("%s took %s", name, elapsed)
}

func RunEmulation() {
	// defer timeTrack(time.Now(), "RunEmulation")
	for {
		cpuTurn = vic.Run(!run)
		if cpu.State == mos6510.ReadInstruction && !run {
			execInst.Lock()
		}
		if cpuTurn {
			cpu.NextCycle()
			if cpu.State == mos6510.ReadInstruction {
				outputDriver.DumpCode(cpu.FullInst)
				if conf.Breakpoint == cpu.InstStart {
					conf.Disassamble = true
					run = false
				}
			}
		}
		cia1.Run()
		cia2.Run()

		if cpu.State == mos6510.ReadInstruction {
			if !run || conf.Disassamble {
				Disassamble()
			}
		}
	}
}

func main() {
	// var exit chan bool
	// exit = make(chan bool)

	confload.Load("config.ini", conf)

	clog.LogLevel = conf.LogLevel
	clog.StartLogging = conf.StartLogging
	if conf.FileLog != "" {
		clog.EnableFileLog(conf.FileLog)
	}

	// f, err := os.Create("newC64.prof")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	setup()
	go input()

	run = true
	cpuTurn = true
	// go func() {
	go RunEmulation()

	// }()

	outputDriver.Run()

	// cpu.DumpStats()
	// <-exit
}

package main

import (
	"fmt"
	"log"
	"newC64/cia6526"
	"newC64/clog"
	"newC64/confload"
	"newC64/graphic"
	"newC64/memory"
	"newC64/mos6510"
	"newC64/pla906114"
	"newC64/vic6569"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/mattn/go-tty"
)

const (
	ramSize     = 65536
	kernalSize  = 8192
	basicSize   = 8192
	ioSize      = 4096
	chargenSize = 4096

	Stopped = 0
	Paused  = 1
	Running = 2
)

var (
	conf = &confload.ConfigData{}

	cpu  mos6510.CPU
	pla  pla906114.PLA
	cia1 cia6526.CIA
	cia2 cia6526.CIA

	mem     memory.MEM
	kernal  memory.MEM
	basic   memory.MEM
	chargen memory.MEM
	io      memory.MEM
	vic     vic6569.VIC
	cycles  int32

	outputDriver graphic.Driver
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
	mem.Init(ramSize, "")
	mem.Clear(true)
	io.Init(ioSize, "")
	io.Clear(false)
	kernal.Init(kernalSize, "assets/roms/kernal.bin")
	basic.Init(basicSize, "assets/roms/basic.bin")
	chargen.Init(chargenSize, "assets/roms/char.bin")

	// PLA Setup
	pla.Init(&mem.Val[1], conf)
	pla.Attach(&mem, pla906114.RAM, 0)
	pla.Attach(&io, pla906114.IO, pla906114.IOStart)
	pla.Attach(&kernal, pla906114.KERNAL, pla906114.KernalStart)
	pla.Attach(&basic, pla906114.BASIC, pla906114.BasicStart)
	pla.Attach(&chargen, pla906114.CHAR, pla906114.CharStart)

	outputDriver = &graphic.SDLDriver{}
	vic.Init(&mem, &io, &chargen, outputDriver, conf)

	// CPU Setup
	cpu.Init(&pla, &vic.SystemClock, conf)

	cia1.Init("CIA1", io.GetView(0x0C00, 0x0200), &vic.SystemClock)
	outputDriver.SetKeyboardLine(&cia1.InputLine)
	cia2.Init("CIA2", io.GetView(0x0D00, 0x0200), &vic.SystemClock)

	pla.Connect(&vic, &cia1, &cia2)

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
			cia1.Dump()
		case 's':
			Disassamble()
			pla.DumpStack(cpu.SP)
		case 'z':
			Disassamble()
			pla.Dump(0)
		case 'r':
			conf.Disassamble = false
			run = true
			execInst.Unlock()
		case 'l':
			// LoadPRG(&pla, "./prg/GARDEN.prg")
			LoadPRG(&pla, conf.LoadPRG)
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
		case 'q':
			cpu.DumpStats()
			os.Exit(0)
		default:
			dumpAddr += string(r)
			fmt.Printf("%c", r)
			if len(dumpAddr) == 4 {
				hx, _ := strconv.ParseInt(dumpAddr, 16, 64)
				pla.Dump(uint16(hx))
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
	fmt.Printf("%s\n", cpu.Disassemble())
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Now().Sub(start)
	log.Printf("%s took %s", name, elapsed)
}

func RunEmulation() {
	// defer timeTrack(time.Now(), "RunEmulation")
	if cpu.State == mos6510.ReadInstruction && !run {
		execInst.Lock()
	}

	cpuTurn = vic.Run(!run)
	if cpuTurn {
		cpu.NextCycle()
		if cpu.State == mos6510.ReadInstruction {
			if conf.Breakpoint == cpu.InstStart {
				conf.Disassamble = true
				run = false
			}
		}
	}
	cia1.Run()
	cia2.Run()

	// for i := 0; i < 1000; i++ {
	// 	// pause
	// }

	if cpu.State == mos6510.ReadInstruction && !run {
		Disassamble()
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
	for {
		RunEmulation()
	}
	// }()

	// outputDriver.Run()

	// cpu.DumpStats()
	// <-exit
}

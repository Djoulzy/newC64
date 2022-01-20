package main

import (
	"fmt"
	"newC64/cia6526"
	"newC64/clog"
	"newC64/confload"
	"newC64/graphic"
	"newC64/memory"
	"newC64/mos6510"
	"newC64/pla906114"
	"newC64/vic6569"
	"os"
	"time"
)

var conf = &confload.ConfigData{}

const (
	ramSize     = 65536
	kernalSize  = 8192
	basicSize   = 8192
	ioSize      = 4096
	chargenSize = 4096
)

var (
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
	exitProcess  chan bool
	cmd          chan rune
)

// func init() {
// 	// This is needed to arrange that main() runs on main thread.
// 	// See documentation for functions that are only allowed to be called from the main thread.
// 	runtime.LockOSThread()
// }

func setup() {
	// ROMs & RAM Setup
	mem.Init(ramSize, "")
	io.Init(ioSize, "")
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

	outputDriver = &graphic.GIODriver{}
	vic.Init(&mem, &io, &chargen, outputDriver, conf)

	// CPU Setup
	cpu.Init(&pla, &vic.SystemClock, conf)

	cia1.Init("CIA1", io.GetView(0x0C00, 0x0200), &vic.SystemClock)
	cia2.Init("CIA2", io.GetView(0x0D00, 0x0200), &vic.SystemClock)

	pla.Connect(&vic, &cia1, &cia2)

	vic.IRQ_Pin = &cpu.IRQ_pin
	cia1.Signal_Pin = &cpu.IRQ_pin
	cia2.Signal_Pin = &cpu.NMI_pin
}

func Disassamble() {
	fmt.Printf("\n%s %s", vic.Disassemble(), cpu.Disassemble())
}

func RunEmulation() {
	var cpuTurn bool
	time.Sleep(time.Second * 2)
	for {
		cpuTurn = vic.Run()
		if cpuTurn {
			if cpu.State == mos6510.ReadInstruction {
				if cpu.NMI_pin > 0 {
					cpu.NMI()
				}
				if (cpu.IRQ_pin > 0) && (cpu.S & ^mos6510.I_mask) == 0 {
					cpu.IRQ()
				}
			}
			cpu.NextCycle()
		}
		cia1.Run()
		cia2.Run()
	}
}

func main() {
	args := os.Args
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

	if len(args) > 1 {
		// addr, _ := LoadPRG(mem, args[1])
		// cpu.GoTo(addr)
	}

	go RunEmulation()

	outputDriver.Start()
}

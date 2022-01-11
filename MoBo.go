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
	"runtime"
	"strconv"

	"github.com/mattn/go-tty"
)

var conf = &confload.ConfigData{}
var run bool
var step bool

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

	video       graphic.Driver
	exitProcess chan bool
	cmd         chan rune
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func setup() {
	// ROMs & RAM Setup
	mem.Init(ramSize, "")
	io.Init(ioSize, "")
	kernal.Init(kernalSize, "assets/roms/kernal.bin")
	basic.Init(basicSize, "assets/roms/basic.bin")
	chargen.Init(chargenSize, "assets/roms/char.bin")

	// PLA Setup
	pla.Init(&mem.Val[1])
	pla.Attach(&mem, pla906114.RAM, 0)
	pla.Attach(&io, pla906114.IO, pla906114.IOStart)
	pla.Attach(&kernal, pla906114.KERNAL, pla906114.KernalStart)
	pla.Attach(&basic, pla906114.BASIC, pla906114.BasicStart)
	pla.Attach(&chargen, pla906114.CHAR, pla906114.CharStart)

	if conf.Display {
		video = &graphic.SDLDriver{}
		vic.Init(&mem, &io, &chargen, video)
	} else {
		vic.SystemClock = 0
	}

	// CPU Setup
	cpu.Init(&pla, &vic.SystemClock, conf)

	cia1.Init("CIA1", io.GetView(0x0C00, 0x0200), &vic.SystemClock)
	cia2.Init("CIA2", io.GetView(0x0D00, 0x0200), &vic.SystemClock)

	vic.IRQ_Pin = &cpu.IRQ
	cia1.Signal_Pin = &cpu.IRQ
	cia2.Signal_Pin = &cpu.NMI
}

func input() {
	exitProcess = make(chan bool)
	cmd = make(chan rune)
	var keyb *tty.TTY
	keyb, _ = tty.Open()
	for {
		r, _ := keyb.ReadRune()
		cmd <- r
	}
}

func main() {
	var cpuTurn bool

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
	pla.DumpChar(0xD8)

	if len(args) > 1 {
		// addr, _ := LoadPRG(mem, args[1])
		// cpu.GoTo(addr)
	}

	run = true
	step = false
	cpuTurn = true
	dumpAddr := ""
	go input()

ENDPROCESS:
	for {
		select {
		case ch := <-cmd:
			switch ch {
			case 's':
				cpu.Disassemble()
				pla.DumpStack(cpu.SP)
			case 'z':
				cpu.Disassemble()
				pla.Dump(0)
			case 'c':
				run = true
				step = false
			case ' ':
				step = true
				run = !run
				fmt.Printf("\n(s) Stack Dump - (z) Zero Page - (c) Continue - (sp) Pause / unpause > ")
			case 'q':
				break ENDPROCESS
			default:
				dumpAddr += string(ch)
				fmt.Printf("%c", ch)
				if len(dumpAddr) == 4 {
					hx, _ := strconv.ParseInt(dumpAddr, 16, 64)
					pla.Dump(uint16(hx))
					dumpAddr = ""
				}
			}
		default:
			if run {
				if conf.Display {
					cpuTurn = vic.Run()
				} else {
					vic.SystemClock++
				}
				if cpuTurn {
					cpu.NextCycle()
				}
				cia1.Run()
				cia2.Run()
			}
			if step {
				if cpu.State == mos6510.ReadInstruction {
					run = false
				}
			}
			if conf.Breakpoint == cpu.InstStart && cpu.State == mos6510.ReadInstruction {
				run = false
				step = true
			}
		}
	}
}

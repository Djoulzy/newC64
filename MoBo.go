package main

import (
	"log"
	"newC64/clog"
	"newC64/confload"
	"newC64/graphic"
	"newC64/mos6510"
	"newC64/pla906114"
	"newC64/vic6569"
	"os"
	"runtime"

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
	cpu mos6510.CPU
	pla pla906114.PLA

	mem     []byte
	kernal  []byte
	basic   []byte
	chargen []byte
	io      []byte
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
	mem = make([]byte, ramSize)
	io = make([]byte, ioSize)
	kernal = make([]byte, kernalSize)
	basic = make([]byte, basicSize)
	chargen = make([]byte, chargenSize)

	// PLA Setup
	pla.Init()
	pla.Attach(mem, pla906114.RAM, 0)
	pla.Attach(io, pla906114.IO, pla906114.IOStart)
	pla.Attach(kernal, pla906114.KERNAL, pla906114.KernalStart)
	pla.Attach(basic, pla906114.BASIC, pla906114.BasicStart)
	pla.Attach(chargen, pla906114.CHAR, pla906114.CharStart)

	pla.Load(pla906114.KERNAL, "assets/roms/kernal.bin")
	pla.Load(pla906114.BASIC, "assets/roms/basic.bin")
	pla.Load(pla906114.CHAR, "assets/roms/char.bin")

	// CPU Setup
	cpu.Init(&pla, conf)

	if conf.Display {
		video = &graphic.SDLDriver{}
		vic.Init(mem, io, chargen, video)
	}
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
		addr, _ := LoadPRG(mem, args[1])
		cpu.GoTo(addr)
	}

	run = true
	step = false
	go input()

ENDPROCESS:
	for {
		select {
		case ch := <-cmd:
			switch ch {
			case 's':
				cpu.Disassemble()
				pla.DumpStack(cpu.SP)
			case 'd':
				cpu.Disassemble()
				pla.Dump(conf.Dump)
			case 'z':
				cpu.Disassemble()
				pla.Dump(0)
			case 'c':
				run = true
				step = false
			case ' ':
				step = true
				run = !run
			case 'q':
				break ENDPROCESS
			case 't':
				log.Println("test")
			}
		default:
			if run {
				if conf.Display {
					vic.Run()
				}
				cpu.NextCycle()
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

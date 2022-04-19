package pla906114

import (
	"fmt"
	"log"
	"newC64/cia6526"
	"newC64/clog"
	"newC64/confload"
	"newC64/memory"
	"newC64/trace"
	"newC64/vic6569"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (P *PLA) Init(settings *byte, conf *confload.ConfigData) {
	P.setting = settings
	P.conf = conf
	P.count = 0
}

func (P *PLA) Attach(mem *memory.MEM, memtype MemType, startLocation int) {
	P.Mem[memtype] = mem
	P.startLocation[memtype] = startLocation
}

func (P *PLA) Connect(vic *vic6569.VIC, cia1 *cia6526.CIA, cia2 *cia6526.CIA) {
	P.vic = vic
	P.cia1 = cia1
	P.cia2 = cia2
}

func (P *PLA) Clear(memtype MemType) {
	P.Mem[memtype].Clear(false)
}

func (P *PLA) getChip(addr uint16) MemType {
	if addr < BasicStart { // Premiere Zone de RAM: 0000 -> A000
		// if addr == 0xF5 && P.count == 1 {
		// 	os.Exit(1)
		// } else {
		// 	P.count++
		// }
		return RAM
	}
	if addr < BasicEnd {
		if *P.setting&3 == 3 {
			return BASIC
		} else {
			return RAM
		}
	}
	if addr < IOStart {
		return RAM
	}
	if addr < KernalStart {
		if *P.setting&(HIRAM|LORAM) == 0 {
			return RAM
		}
		if *P.setting&CHAREN == 0 {
			return CHAR
		}
		if *P.setting&CHAREN == CHAREN {
			return IO
		}
		log.Fatal("Bad memory zone")
	}
	if addr > CharEnd {
		if *P.setting&HIRAM == HIRAM {
			return KERNAL
		} else {
			return RAM
		}
	}
	log.Fatal("Bad memory zone")
	return RAM
}

func (P *PLA) Read(addr uint16) byte {
	// defer timeTrack(time.Now(), "Read")
	dest := P.getChip(addr)
	transAddr := addr - uint16(P.startLocation[dest])

	if dest == IO {
		if transAddr < 0x0400 {
			return P.vic.Read(transAddr)
		}
		if transAddr < 0x0800 {
			// log.Fatal("SID Not implemented")
			return P.Mem[IO].Val[transAddr]
		}
		if transAddr < 0x0C00 {
			return P.Mem[IO].Val[transAddr]
		}
		if transAddr < 0x0D00 {
			return P.cia1.Read(transAddr - 0x0C00)
		}
		if transAddr < 0x0E00 {
			return P.cia2.Read(transAddr - 0x0D00)
		} else {
			// log.Fatal("I/O Not implemented")
			return P.Mem[IO].Val[transAddr]
		}
	}

	// fmt.Printf("pla906114 - Read - %04X - Zone: %d\n", addr, dest)
	return P.Mem[dest].Val[transAddr]
}

func (P *PLA) Write(addr uint16, value byte) {
	var transAddr uint16
	// if addr >= 0xE000 {
	// 	log.Printf("Write to %04X", addr)
	// 	os.Exit(1)
	// }
	if P.getChip(addr) == IO {
		transAddr = addr - uint16(P.startLocation[IO])
		if transAddr < 0x0400 {
			P.vic.Write(transAddr, value)
			return
		}
		if transAddr < 0x0800 {
			// log.Fatal("SID Not implemented")
			P.Mem[IO].Val[transAddr] = value
			// P.Mem[RAM].Val[addr] = value
			return
		}
		if transAddr < 0x0C00 {
			P.Mem[IO].Val[transAddr] = value
			// P.Mem[RAM].Val[addr] = value
			return
		}
		if transAddr < 0x0D00 {
			P.cia1.Write(transAddr-0x0C00, value)
			return
		}
		if transAddr < 0x0E00 {
			P.cia2.Write(transAddr-0x0D00, value)
			return
		} else {
			// log.Fatal("I/O Not implemented")
			P.Mem[IO].Val[transAddr] = value
			// P.Mem[RAM].Val[addr] = value
			return
		}
	}
	P.Mem[RAM].Val[addr] = value
}

func (P *PLA) GetView(start int, size int) []byte {
	return P.Mem[RAM].Val[start : start+size]
}

func (P *PLA) Dump(startAddr uint16) {
	var val byte
	var line string
	var ascii string

	cpt := startAddr
	fmt.Printf("\n")
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		line = ""
		ascii = ""
		for i := 0; i < 16; i++ {
			// val = P.Read(cpt)
			val = P.Mem[RAM].Val[cpt]
			if val != 0x00 && val != 0xFF {
				line = line + clog.CSprintf("white", "black", "%02X", val) + " "
			} else {
				line = fmt.Sprintf("%s%02X ", line, val)
			}
			if _, ok := trace.PETSCII[val]; ok {
				ascii += fmt.Sprintf("%s", string(trace.PETSCII[val]))
			} else {
				ascii += "."
			}
			cpt++
		}
		fmt.Printf("%s - %s\n", line, ascii)
	}
}

func (P *PLA) DumpStack(sp byte) {
	cpt := uint16(0x0100)
	fmt.Printf("\n")
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		for i := 0; i < 16; i++ {
			if cpt == StackStart+uint16(sp) {
				clog.CPrintf("white", "red", "%02X", P.Mem[RAM].Val[cpt])
				fmt.Print(" ")
				// fmt.Printf("%c[41m%c[0m[0;31m%02X%c[0m ", 27, 27, P.Mem[RAM].Val[cpt], 27)
			} else {
				fmt.Printf("%02X ", P.Mem[RAM].Val[cpt])
			}
			cpt++
		}
		fmt.Println()
	}
}

func (P *PLA) DumpChar(screenCode byte) {
	cpt := uint16(screenCode) << 3
	for j := 0; j < 4; j++ {
		for i := 0; i < 8; i++ {
			fmt.Printf("%04X : %08b\n", cpt, P.Mem[CHAR].Val[cpt])
			cpt++
		}
		fmt.Println()
	}
}

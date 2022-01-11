package pla906114

import (
	"fmt"
	"log"
	"newC64/clog"
	"newC64/memory"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (P *PLA) Init(settings *byte) {
	P.setting = settings
}

func (P *PLA) Attach(mem *memory.MEM, memtype MemType, startLocation int) {
	P.Mem[memtype] = mem
	P.startLocation[memtype] = startLocation
}

func (P *PLA) Clear(memtype MemType) {
	P.Mem[memtype].Clear()
}

func (P *PLA) getChip(addr uint16) MemType {
	if addr < BasicStart { // Premiere Zone de RAM: 0000 -> A000
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
	destAddr := addr - uint16(P.startLocation[dest])
	// fmt.Printf("pla906114 - Read - %04X - Zone: %d\n", addr, dest)
	return P.Mem[dest].Val[destAddr]
}

func (P *PLA) Write(addr uint16, value byte) {
	var transAddr uint16
	// if addr > 0x0400 && addr < 0x07E7 {
	// 	os.Exit(1)
	// }
	if P.getChip(addr) == IO {
		transAddr = addr - uint16(P.startLocation[IO])
		if addr < 0xD400 {
			P.Mem[IO].VicRegWrite(transAddr, value)
			return
		}
		if addr < 0xDC00 {
			P.Mem[IO].Val[transAddr] = value
			P.Mem[IO].Written[transAddr] = true
			return
		}
		if addr < 0xDE00 {
			P.Mem[IO].CiaRegWrite(transAddr, value)
			return
		} else {
			log.Fatal("Bad IO addr")
		}
	}
	P.Mem[RAM].Val[addr] = value
}

func (P *PLA) GetView(start int, size int) []byte {
	return P.Mem[RAM].Val[start : start+size]
}

func (P *PLA) Dump(startAddr uint16) {
	var val byte
	var zone MemType
	var transAddr uint16
	cpt := startAddr
	fmt.Printf("\n")
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		for i := 0; i < 16; i++ {
			zone = P.getChip(cpt)
			transAddr = cpt - uint16(P.startLocation[zone])
			val = P.Mem[zone].Val[transAddr]
			if P.Mem[zone].Written[transAddr] {
				clog.CPrintf("white", "red", "%02X", val)
				fmt.Print(" ")
			} else if val != 0x00 && val != 0xFF {
				clog.CPrintf("white", "blue", "%02X", val)
				fmt.Print(" ")
			} else {
				fmt.Printf("%02X ", val)
			}
			cpt++
		}
		fmt.Println()
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

package pla906114

import (
	"fmt"
	"log"
	"newC64/memory"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func (P *PLA) Init() {
	P.setting = 7
}

func (P *PLA) Attach(mem *memory.MEM, memtype MemType, startLocation int) {
	P.Mem[memtype] = mem
	P.startLocation[memtype] = startLocation
}

func (P *PLA) Clear(memtype MemType) {
	P.Mem[memtype].Clear()
}

// func (P *PLA) Load(memtype MemType, filename string) {

// 	P.Mem[memtype].readOnly = true

// 	data, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if len(data) != P.Mem[memtype].Size {
// 		panic("Bad ROM Size")
// 	}
// 	for i := 0; i < P.Mem[memtype].Size; i++ {
// 		P.Mem[memtype].Cells[i] = byte(data[i])
// 	}
// }

func (P *PLA) getChip(addr uint16) MemType {
	switch {
	case addr < BasicStart:
		return RAM
	case addr < BasicEnd:
		if P.setting&3 == 3 {
			return BASIC
		} else {
			return RAM
		}
	case addr < CharStart:
		return RAM
	case addr < KernalStart:
		if P.setting&3 == 0 {
			return RAM
		} else if P.setting&4 == 0 {
			return CHAR
		} else if P.setting == 1 {
			return RAM
		} else {
			return IO
		}
	default:
		if P.setting&3 < 2 {
			return RAM
		} else {
			return KERNAL
		}
	}
}

// func (P *PLA) getChip(addr uint16) MemType {
// 	if addr >= KernalStart {
// 		return KERNAL
// 	}
// 	if addr >= IOStart {
// 		return IO
// 	}
// 	if addr > BasicEnd {
// 		return RAM
// 	}
// 	if addr >= BasicStart {
// 		return BASIC
// 	}

// 	return RAM
// }

func (P *PLA) Read(addr uint16) byte {
	// defer timeTrack(time.Now(), "Read")
	dest := P.getChip(addr)
	destAddr := addr - uint16(P.startLocation[dest])
	// fmt.Printf("pla906114 - Read - %04X - Zone: %d\n", addr, dest)
	return P.Mem[dest].Val[destAddr]
}

func (P *PLA) Write(addr uint16, value byte) {
	if P.getChip(addr) == IO {
		if addr < 0xD400 {
			P.Mem[IO].VicRegWrite(addr-uint16(P.startLocation[IO]), value)
			return
		}
		if addr < 0xDC00 {
			P.Mem[IO].Val[addr-uint16(P.startLocation[IO])] = value
			return
		}
		if addr < 0xDE00 {
			P.Mem[IO].CiaRegWrite(addr-uint16(P.startLocation[IO]), value)
			return
		}
	}
	P.Mem[RAM].Val[addr] = value
}

func (P *PLA) GetView(start int, size int) []byte {
	return P.Mem[RAM].Val[start : start+size]
}

func (P *PLA) Dump(startAddr uint16) {
	cpt := startAddr
	fmt.Printf("\n")
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		for i := 0; i < 16; i++ {
			fmt.Printf("%02X ", P.Read(cpt))
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
				fmt.Printf("%c[0;31m%02X%c[0m ", 27, P.Mem[RAM].Val[cpt], 27)
			} else {
				fmt.Printf("%02X ", P.Mem[RAM].Val[cpt])
			}
			cpt++
		}
		fmt.Println()
	}
}

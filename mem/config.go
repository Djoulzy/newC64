package mem

import (
	"fmt"
	"io/ioutil"
	"newC64/clog"
	"newC64/trace"
)

const (
	PAGE_DIVIDER = 12
	READWRITE    = false
	READONLY     = true
)

type MEMAccess interface {
	MRead([]byte, uint16) byte
	MWrite([]byte, uint16, byte)
}

type CONFIG struct {
	Layers       [][]byte    // Liste des couches de memoire
	LayersName   []string    // Nom de la couche
	Start        []uint16    // Addresse de début de la couche
	PagesUsed    [][]bool    // Pages utilisées par la couche
	ReadOnly     []bool      // Mode d'accès à la couche
	LayerByPages []int       // Couche active pour la page
	Accessors    []MEMAccess // Reader/Writer de la couche
	TotalPages   int         // Nb total de pages
}

func InitConfig(nbLayers int, size int) CONFIG {
	C := CONFIG{}
	C.Layers = make([][]byte, nbLayers)
	C.LayersName = make([]string, nbLayers)
	C.Start = make([]uint16, nbLayers)
	C.TotalPages = int(size >> PAGE_DIVIDER)
	C.LayerByPages = make([]int, C.TotalPages)
	C.PagesUsed = make([][]bool, nbLayers)
	C.ReadOnly = make([]bool, nbLayers)
	C.Accessors = make([]MEMAccess, nbLayers)
	return C
}

func LoadROM(size int, file string) []byte {
	val := make([]byte, size)
	if len(file) > 0 {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		if len(data) != size {
			panic("Bad ROM Size")
		}
		for i := 0; i < size; i++ {
			val[i] = byte(data[i])
		}
	}
	return val
}

func Clear(zone []byte) {
	cpt := 0
	fill := byte(0x00)
	for i := range zone {
		zone[i] = fill
		cpt++
		if cpt == 0x40 {
			fill = ^fill
			cpt = 0
		}
	}
}

func (C *CONFIG) Attach(name string, layerNum int, pageNum int, content []byte, mode bool) {
	nbPages := len(content) >> PAGE_DIVIDER
	C.LayersName[layerNum] = name
	C.Layers[layerNum] = content
	C.Start[layerNum] = uint16(pageNum << PAGE_DIVIDER)
	C.ReadOnly[layerNum] = mode
	C.PagesUsed[layerNum] = make([]bool, C.TotalPages)
	for i := 0; i < C.TotalPages; i++ {
		C.PagesUsed[layerNum][i] = false
	}
	for i := 0; i < nbPages; i++ {
		C.LayerByPages[pageNum+i] = layerNum
		C.PagesUsed[layerNum][pageNum+i] = true
	}
	C.Accessors[layerNum] = C
}

func (C *CONFIG) Accessor(layerNum int, access MEMAccess) {
	C.Accessors[layerNum] = access
}

func (C *CONFIG) MRead(mem []byte, addr uint16) byte {
	// clog.Test("MEM", "MRead", "Addr: %04X -> %02X", addr, mem[addr])
	return mem[addr]
}

func (C *CONFIG) MWrite(mem []byte, addr uint16, val byte) {
	// clog.Test("MEM", "MWrite", "Addr: %04X -> %02X", addr, val)
	mem[addr] = val
}

type BANK struct {
	Selector *byte
	Layouts  []CONFIG
}

func InitBanks(nbMemLayout int, sel *byte) BANK {
	B := BANK{}
	B.Layouts = make([]CONFIG, nbMemLayout)
	B.Selector = sel
	return B
}

func (B *BANK) Read(addr uint16) byte {
	// clog.Test("MEM", "Read", "Addr: %04X, Page: %d, Selector: %d", addr, int(addr>>PAGE_DIVIDER), *B.Selector&0x1F)
	bank := B.Layouts[*B.Selector&0x1F]
	layerNum := bank.LayerByPages[int(addr>>PAGE_DIVIDER)]
	// return C.Layers[layerNum][addr-C.Start[layerNum]]
	// clog.Test("MEM", "Read", "Addr: %04X, Page: %d, Layer: %d", addr, int(addr>>PAGE_DIVIDER), layerNum)
	return bank.Accessors[layerNum].MRead(bank.Layers[layerNum], addr-bank.Start[layerNum])
}

func (B *BANK) Write(addr uint16, value byte) {
	bank := B.Layouts[*B.Selector&0x1F]
	layerNum := bank.LayerByPages[int(addr>>PAGE_DIVIDER)]
	if bank.ReadOnly[layerNum] {
		layerNum = 0
	}
	// clog.Test("MEM", "Write", "Addr: %04X, Page: %d, Layer: %d", addr, int(addr>>PAGE_DIVIDER), layerNum)
	bank.Accessors[layerNum].MWrite(bank.Layers[layerNum], addr-bank.Start[layerNum], value)
}

func (C *CONFIG) Show() {
	clog.CPrintf("dark_gray", "black", "\n%10s: ", "Pages")
	for p := range C.LayerByPages {
		clog.CPrintf("dark_gray", "black", " %02d  ", p)
	}
	clog.CPrintf("dark_gray", "black", "\n%10s: ", "Start Addr")
	for p := range C.LayerByPages {
		clog.CPrintf("light_gray", "black", "%04X ", p<<PAGE_DIVIDER)
	}
	fmt.Printf("\n")
	for layerRead := range C.Layers {
		clog.CPrintf("light_gray", "black", "%10s: ", C.LayersName[layerRead])
		for pagenum, layerFound := range C.LayerByPages {
			if C.PagesUsed[layerRead][pagenum] {
				if layerFound == layerRead {
					if C.ReadOnly[layerRead] {
						clog.CPrintf("black", "yellow", "     ")
					} else {
						clog.CPrintf("black", "green", "     ")
					}
				} else {
					if C.ReadOnly[layerFound] && !C.ReadOnly[layerRead] {
						clog.CPrintf("black", "red", "     ")
					} else {
						clog.CPrintf("black", "light_gray", "     ")
					}
				}
			} else {
				clog.CPrintf("black", "dark_gray", "     ")
			}
		}
		fmt.Printf(" - %d\n", layerRead)
	}
	clog.CPrintf("dark_gray", "black", "\n%12s", " ")
	clog.CPrintf("black", "green", "%s", " Read/Write ")
	clog.CPrintf("black", "black", "%s", "  ")
	clog.CPrintf("black", "yellow", "%s", " Read Only ")
	clog.CPrintf("black", "black", "%s", "  ")
	clog.CPrintf("black", "red", "%s", " Write Only ")
	clog.CPrintf("black", "black", "%s", "  ")
	clog.CPrintf("black", "light_gray", "%s", " Masked ")
	clog.CPrintf("black", "black", "%s", " ")
	fmt.Printf("\n\n")
}

func (B *BANK) Dump(startAddr uint16) {
	var val byte
	var line string
	var ascii string

	cpt := startAddr
	for j := 0; j < 16; j++ {
		fmt.Printf("%04X : ", cpt)
		line = ""
		ascii = ""
		for i := 0; i < 16; i++ {
			val = B.Read(cpt)
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

func (B *BANK) Show() {
	B.Layouts[*B.Selector&0x1F].Show()
}

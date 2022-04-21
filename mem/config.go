package mem

import (
	"fmt"
	"newC64/clog"
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

func (C *CONFIG) MRead(mem []byte, translatedAddr uint16) byte {
	// clog.Test("MEM", "MRead", "Addr: %04X -> %02X", addr, mem[addr])
	return mem[translatedAddr]
}

func (C *CONFIG) MWrite(mem []byte, translatedAddr uint16, val byte) {
	// clog.Test("MEM", "MWrite", "Addr: %04X -> %02X", addr, val)
	mem[translatedAddr] = val
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
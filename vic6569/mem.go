package vic6569

import (
	"fmt"
	"newC64/clog"
)

const PAGE_DIVIDER = 12

type VicMem struct {
	Pages      []int
	Layers     [][]byte
	LayersName []string
}

func (VM *VicMem) Init(nbLayers int, size uint16) *VicMem {
	VM.Layers = make([][]byte, nbLayers)
	VM.LayersName = make([]string, nbLayers)
	nbPages := int(size >> PAGE_DIVIDER)
	VM.Pages = make([]int, nbPages)
	return VM
}

func (VM *VicMem) Attach(name string, layerNum int, pageNum int, content []byte) {
	nbPages := len(content) >> PAGE_DIVIDER
	VM.LayersName[layerNum] = name
	VM.Layers[layerNum] = content
	for i := 0; i < nbPages; i++ {
		VM.Pages[pageNum+i] = layerNum
	}
}

func (VM *VicMem) Read(addr uint16) byte {
	page := int(addr >> PAGE_DIVIDER)
	pageStart := uint16(page << PAGE_DIVIDER)
	return VM.Layers[VM.Pages[page]][addr-pageStart]
}

func (VM *VicMem) Show() {
	clog.CPrintf("darkgray", "black", "%10s: ", " ")
	for p := range VM.Pages {
		clog.CPrintf("darkgray", "black", "%02d ", p)
	}
	fmt.Printf("\n")
	for l := range VM.Layers {
		clog.CPrintf("darkgray", "black", "%10s: ", VM.LayersName[l])
		for _, page := range VM.Pages {
			if page == l {
				clog.CPrintf("white", "black", " X ")
			} else {
				clog.CPrintf("darkgray", "black", " - ")
			}
		}
		fmt.Printf("\n")
	}
}

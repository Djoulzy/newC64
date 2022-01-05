package vic6569

import (
	"newC64/graphic"
	"newC64/memory"
)

var (
	Black      byte = 0
	White      byte = 1
	Red        byte = 2
	Cyan       byte = 3
	Violet     byte = 4
	Green      byte = 5
	Blue       byte = 6
	Yellow     byte = 7
	Orange     byte = 8
	Brown      byte = 9
	Lightred   byte = 10
	Darkgrey   byte = 11
	Grey       byte = 12
	Lightgreen byte = 13
	Lightblue  byte = 14
	Lightgrey  byte = 15
)

var Colors [16]graphic.RGB = [16]graphic.RGB{
	{R: 0, G: 0, B: 0},       // Black
	{R: 255, G: 255, B: 255}, // White
	{R: 137, G: 78, B: 67},   // Red
	{R: 146, G: 195, B: 203}, // Cyan
	{R: 138, G: 87, B: 176},  // Violet
	{R: 128, G: 174, B: 89},  // Green
	{R: 68, G: 63, B: 164},   // Blue
	{R: 215, G: 221, B: 137}, // Yellow
	{R: 146, G: 106, B: 56},  // Orange
	{R: 100, G: 82, B: 23},   // Brown
	{R: 184, G: 132, B: 122}, // Lightred
	{R: 96, G: 96, B: 96},    // Darkgrey
	{R: 138, G: 138, B: 138}, // Grey
	{R: 191, G: 233, B: 155}, // Lightgreen
	{R: 131, G: 125, B: 216}, // Lightblue
	{R: 179, G: 179, B: 179}, // Lightgrey
}

// VIC :
type VIC struct {
	VML    [40]uint16 // Video Matrix Line
	VMLI   byte       // Video Matrix Line Indexer
	VC     uint16     // Vide Counter
	VCBASE uint16     // Video Counter Base
	RC     byte       // Row counter
	BA     bool       // High: normal / Low: BadLine

	beamX int
	beamY int
	cycle int

	visibleArea bool
	displayArea bool
	drawArea    bool

	ColorBuffer [40]byte
	CharBuffer  [40]byte

	IRQ_Pin   *int
	RasterIRQ uint16
	graph     graphic.Driver

	chargen *memory.MEM
	io      *memory.MEM
	color   *memory.MEM
	ram     *memory.MEM
	screen  *memory.MEM
}

const (
	CharStart   = 0xD000
	IOStart     = 0xD000
	colorStart  = 0xD800
	screenStart = 0x0400

	REG_CTRL1  uint16 = 0xD011 // Screen control (0b01111111)
	REG_RASTER uint16 = 0xD012 // Raster 8 first bits
	REG_CTRL2  uint16 = 0xD016 // Screen control (0b01111111)
	REG_IRQ    uint16 = 0xD019 // IRQ Register
	REG_SETIRQ uint16 = 0xD01A // IRQ Enabler
	REG_EC     uint16 = 0xD020 // Border Color
	REG_B0C    uint16 = 0xD021 // Background color 0
	PALNTSC    uint16 = 0x02A6

	YSCROLL byte = 0b00000111 // From REG_CTRL1
	RSEL    byte = 0b00001000 // rom REG_CTRL1 : 0 = 24 rows; 1 = 25 rows.
	DEN     byte = 0b00010000 // rom REG_CTRL1 : 0 = Screen off, 1 = Screen on.
	BMM     byte = 0b00100000 // rom REG_CTRL1 : 0 = Text mode; 1 = Bitmap mode.
	ECM     byte = 0b01000000 // rom REG_CTRL1 : 1 = Extended background mode on.
	RST8    byte = 0b10000000 // rom REG_CTRL1 : Read: Current raster line (bit #8). Write: Raster line to generate interrupt at (bit #8).

	IRQ_RASTER    byte = 0b00000001
	IRQ_SPRT_BG   byte = 0b00000010
	IRQ_SPRT_SPRT byte = 0b00000100
	IRQ_LGTPEN    byte = 0b00001000
)

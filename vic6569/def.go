package vic6569

import (
	"newC64/confload"
	"newC64/graphic"
	"newC64/mem"
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
	VML         [40]uint16 // Video Matrix Line
	VMLI        byte       // Video Matrix Line Indexer
	VC          uint16     // Vide Counter
	VCBASE      uint16     // Video Counter Base
	RC          byte       // Row counter
	BA          bool       // High: normal / Low: BadLine
	SystemClock uint16
	Reg         [64]byte

	conf  *confload.ConfigData
	BeamX int
	BeamY int
	cycle int

	visibleArea bool
	displayArea bool
	drawArea    bool

	ColorBuffer [40]byte
	CharBuffer  [40]byte

	IRQ_Pin   *int
	RasterIRQ uint16
	graph     graphic.Driver

	color      []byte
	BankSel    byte
	ScreenBase uint16
	CharBase   uint16
	bankMem    mem.BANK

	ECM  bool
	BMM  bool
	MCM  bool
	MODE byte
}

const (
	REG_X_SPRT_0           uint16 = iota // X coordinate sprite 0
	REG_Y_SPRT_0                         // Y coordinate sprite 0
	REG_X_SPRT_1                         // X coordinate sprite 1
	REG_Y_SPRT_1                         // Y coordinate sprite 1
	REG_X_SPRT_2                         // X coordinate sprite 2
	REG_Y_SPRT_2                         // Y coordinate sprite 2
	REG_X_SPRT_3                         // X coordinate sprite 3
	REG_Y_SPRT_3                         // Y coordinate sprite 3
	REG_X_SPRT_4                         // X coordinate sprite 4
	REG_Y_SPRT_4                         // Y coordinate sprite 4
	REG_X_SPRT_5                         // X coordinate sprite 5
	REG_Y_SPRT_5                         // Y coordinate sprite 5
	REG_X_SPRT_6                         // X coordinate sprite 6
	REG_Y_SPRT_6                         // Y coordinate sprite 6
	REG_X_SPRT_7                         // X coordinate sprite 7
	REG_Y_SPRT_7                         // Y coordinate sprite 7
	REG_MSBS_X_COOR                      // MSBs of X coordinates
	REG_CTRL1                            // Control register 1
	REG_RASTER                           // Raster counter
	REG_LP_X                             // Light pen X
	REG_LP_Y                             // Light pen Y
	REG_SPRT_ENABLED                     // Sprite enabled
	REG_CTRL2                            // Control register 2
	REG_SPRT_Y_EXP                       // Sprite Y expansion
	REG_MEM_LOC                          // Memory pointers
	REG_IRQ                              // Interrupt register
	REG_IRQ_ENABLED                      // Interrupt enabled
	REG_SPRT_DATA_PRIORITY               // Sprite data priority
	REG_SPRT_MLTCOLOR                    // Sprite multicolor
	REG_SPRT_X_EXP                       // Sprite X expansion
	REG_SPRT_SPRT_COLL                   // Spritesprite collision
	REG_SPRT_DATA_COLL                   // Spritedata collision
	REG_BORDER_COL                       // Border color
	REG_BGCOLOR_0                        // Background color 0
	REG_BGCOLOR_1                        // Background color 1
	REG_BGCOLOR_2                        // Background color 2
	REG_BGCOLOR_3                        // Background color 3
	REG_SPRT_MLTCOLOR_0                  // Sprite multicolor 0
	REG_SPRT_MLTCOLOR_1                  // Sprite multicolor 1
	REG_COLOR_SPRT_0                     // Color sprite 0
	REG_COLOR_SPRT_1                     // Color sprite 1
	REG_COLOR_SPRT_2                     // Color sprite 2
	REG_COLOR_SPRT_3                     // Color sprite 3
	REG_COLOR_SPRT_4                     // Color sprite 4
	REG_COLOR_SPRT_5                     // Color sprite 5
	REG_COLOR_SPRT_6                     // Color sprite 6
	REG_COLOR_SPRT_7                     // Color sprite 7
)

const (
	colorStart  = 0x0800 // 0xD800 translated
	screenStart = 0x0400

	PALNTSC uint16 = 0x02A6

	YSCROLL byte = 0b00000111 // From REG_CTRL1
	RSEL    byte = 0b00001000 // rom REG_CTRL1 : 0 = 24 rows; 1 = 25 rows.
	DEN     byte = 0b00010000 // rom REG_CTRL1 : 0 = Screen off, 1 = Screen on.
	// BMM     byte = 0b00100000 // rom REG_CTRL1 : 0 = Text mode; 1 = Bitmap mode.
	// ECM     byte = 0b01000000 // rom REG_CTRL1 : 1 = Extended background mode on.
	// MCM     byte = 0b00010000 // rom REG_CTRL2
	RST8 byte = 0b10000000 // rom REG_CTRL1 : Read: Current raster line (bit #8). Write: Raster line to generate interrupt at (bit #8).

	IRQ_RST byte = 0b00000001 // Raster line interrupt
	IRQ_MBC byte = 0b00000010 // Sprite collision with background
	IRQ_MMC byte = 0b00000100 // Sprite vs sprite collision
	IRQ_LP  byte = 0b00001000 // Light pen negative edge
)

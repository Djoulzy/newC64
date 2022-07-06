package config

// Globals : Partie globale du fichier de conf
type Globals struct {
	StartLogging bool
	FileLog      string
	RunPerfStats bool
	Disassamble  bool
	LogLevel     int
	Debug        bool
	Display      bool
	LoadPRG      string
	Model        string
	ColorDisplay bool
	CPUModel     string
	Mhz          int
}

type DebugMode struct {
	Breakpoint uint16
	Dump       uint16
	Zone       int
}

// ConfigData : Data structure du fichier de conf
type ConfigData struct {
	Globals
	DebugMode
}

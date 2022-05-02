package config

// Globals : Partie globale du fichier de conf
type Globals struct {
	StartLogging bool
	FileLog      string
	RunPerfStats bool
	Disassamble  bool
	LogLevel     int
	Display      bool
	LoadPRG      string
}

type Debug struct {
	Breakpoint uint16
	Dump       uint16
	Zone       int
}

// ConfigData : Data structure du fichier de conf
type ConfigData struct {
	Globals
	Debug
}

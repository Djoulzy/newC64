package confload

// Globals : Partie globale du fichier de conf
type Globals struct {
	StartLogging bool
	FileLog      string
	Disassamble  bool
	LogLevel     int
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

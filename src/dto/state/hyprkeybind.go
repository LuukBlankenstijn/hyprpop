package state

type HyprKeybind struct {
	Mod        int    `json:"modmask"`
	Key        string `json:"key"`
	Dispatcher string `json:"dispatcher"`
	Arg        string `json:"arg"`
}

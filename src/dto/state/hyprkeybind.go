package state

type HyprKeybind struct {
	Mod        int    `json:"mod"`
	Key        string `json:"key"`
	Dispatcher string `json:"dispatcher"`
	Arg        string `json:"arg"`
}

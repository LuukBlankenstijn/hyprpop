package state

type WindowConfig struct {
	Keybind  Keybind `yaml:"keybind"`
	Location Vec2    `yaml:"position"`
	Name     string  `yaml:"name"`
	Size     Vec2    `yaml:"size"`
	URL      string  `yaml:"url"`
}

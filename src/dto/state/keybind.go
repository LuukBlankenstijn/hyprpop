package state

import (
	"encoding/json"
	"strings"
)

type Keybind struct {
	Mod string `json:"modmask"`
	Key string `json:"key"`
}

const (
	MOD_SHIFT = 1
	MOD_CTRL  = 2
	MOD_CAPS  = 4
	MOD_ALT   = 8
	MOD_NUM   = 16
	MOD_MOD3  = 32
	MOD_SUPER = 64
)

func ModToString(mod int) string {
	switch mod {
	case MOD_SHIFT:
		return "SHIFT"
	case MOD_CTRL:
		return "CTRL"
	case MOD_CAPS:
		return "CAPS"
	case MOD_ALT:
		return "ALT"
	case MOD_NUM:
		return "NUM"
	case MOD_MOD3:
		return "MOD3"
	case MOD_SUPER:
		return "SUPER"
	default:
		return "UNKNOWN"
	}
}

func (k *Keybind) UnmarshalYAML(unmarshal func(any) error) error {
	var keybindStr string
	if err := unmarshal(&keybindStr); err != nil {
		return err
	}

	parts := strings.SplitN(keybindStr, "+", 2)
	if len(parts) == 2 {
		k.Mod = parts[0]
		k.Key = parts[1]
	} else {
		k.Mod = ""
		k.Key = parts[0]
	}

	return nil
}

func (k *Keybind) UnmarshalJSON(data []byte) error {
	var raw struct {
		Mod int    `json:"modmask"`
		Key string `json:"key"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	k.Mod = ModToString(raw.Mod)
	k.Key = raw.Key

	return nil
}

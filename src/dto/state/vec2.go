package state

import (
	"encoding/json"
	"fmt"
)

type Vec2 struct {
	X VectorValue
	Y VectorValue
}

func (v *Vec2) UnmarshalJSON(data []byte) error {
	var arr [2]VectorValue
	if err := json.Unmarshal(data, &arr); err != nil {
		return fmt.Errorf("failed to unmarshal Vec2: %w", err)
	}
	v.X = arr[0]
	v.Y = arr[1]
	return nil
}

func (v *Vec2) UnmarshalYAML(unmarshal func(any) error) error {
	var arr [2]VectorValue
	if err := unmarshal(&arr); err != nil {
		return fmt.Errorf("failed to unmarshal Vec2 from YAML: %w", err)
	}
	v.X = arr[0]
	v.Y = arr[1]
	return nil
}

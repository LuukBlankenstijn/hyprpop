package state

import (
	"encoding/json"
	"fmt"
	"strconv"
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

func (size *Vec2) GetExactSize(monitor Monitor) (string, string, error) {
	var x, y string
	if size.X.IsPercentage {
		x = strconv.Itoa(int(size.X.Value * float64(monitor.GetWidth())))
	} else {
		x = strconv.Itoa(int(size.X.Value))
	}
	if size.Y.IsPercentage {
		y = strconv.Itoa(int(size.Y.Value * float64(monitor.GetHeight())))
	} else {
		y = strconv.Itoa(int(size.Y.Value))
	}
	return x, y, nil
}

func (position *Vec2) GetExactPosition(monitor Monitor) (string, string, error) {
	var pixels int

	// X value
	v := position.X
	if v.IsPercentage {
		pixels = int(v.Value * float64(monitor.GetWidth()))
	} else {
		pixels = int(v.Value)
	}
	x := strconv.Itoa(pixels + monitor.X)

	// Y value
	v = position.Y
	if v.IsPercentage {
		pixels = int(v.Value * float64(monitor.GetHeight()))
	} else {
		pixels = int(v.Value)
	}
	y := strconv.Itoa(pixels + monitor.Y)
	return x, y, nil
}

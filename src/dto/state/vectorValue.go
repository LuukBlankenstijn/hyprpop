package state

import (
	"encoding/json"
	"errors"
	"strconv"
)

type VectorValue struct {
	Value        float64
	IsPercentage bool
}

func (p *VectorValue) UnmarshalYAML(unmarshal func(any) error) error {
	// Try to unmarshal as a float (percentage)
	var floatValue float64
	if err := unmarshal(&floatValue); err == nil {
		if floatValue > 0 && floatValue < 1 {
			p.Value = floatValue
			p.IsPercentage = true
			return nil
		}
	}

	// Try to unmarshal as an int (absolute pixels)
	var intValue int
	if err := unmarshal(&intValue); err == nil {
		p.Value = float64(intValue)
		p.IsPercentage = false
		return nil
	}

	return errors.New("invalid size format: must be an integer (pixels) or a float (percentage)")
}

func (p *VectorValue) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as an integer (pixels)
	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		p.Value = float64(intValue)
		p.IsPercentage = false
		return nil
	}

	// Try to unmarshal as a float (percentage)
	var floatValue float64
	if err := json.Unmarshal(data, &floatValue); err == nil {
		if floatValue < 0 || floatValue > 1 {
			return errors.New("size percentage must be between 0 and 1")
		}
		p.Value = floatValue
		p.IsPercentage = true
		return nil
	}

	return errors.New("invalid size format: must be an integer (pixels) or a float (percentage)")
}

func (p *VectorValue) GetAsString() string {
	if p.IsPercentage {
		return strconv.Itoa(int(p.Value*100)) + "%"
	}
	return strconv.Itoa(int(p.Value))
}

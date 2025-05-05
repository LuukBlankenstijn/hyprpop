package hypr

import (
	"fmt"
	"hyprpop/src/dto/state"
	"strconv"
)

func getExactPosition(p state.Vec2) (string, string, error) {
	monitor, err := getActiveMonitor()
	if err != nil {
		return "", "", fmt.Errorf("could not get current monitor")
	}
	var pixels int

	// X value
	v := p.X
	if v.IsPercentage {
		pixels = int(v.Value * float64(monitor.GetWidth()))
	} else {
		pixels = int(v.Value)
	}
	// fmt.Println(pixels)
	x := strconv.Itoa(pixels + monitor.X)

	// Y value
	v = p.Y
	if v.IsPercentage {
		pixels = int(v.Value * float64(monitor.GetHeight()))
	} else {
		pixels = int(v.Value)
	}
	// fmt.Println(pixels)
	y := strconv.Itoa(pixels + monitor.Y)
	return x, y, nil
}

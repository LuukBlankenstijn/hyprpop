package hypr

import (
	"hyprpop/src/dto/state"
	"strconv"
)

func getExactSize(size state.Vec2, monitorId int) (string, string, error) {
	monitor, err := getMonitorById(monitorId)
	if err != nil {
		return "", "", err
	}
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

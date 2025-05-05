package state

type Monitor struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Workspace Workspace `json:"activeWorkspace"`
	Scale     float64   `json:"scale"`
	Focused   bool      `json:"focused"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
}

func (m *Monitor) GetWidth() int {
	return int(float64(m.Width) / m.Scale)
}

func (m *Monitor) GetHeight() int {
	return int(float64(m.Height) / m.Scale)
}

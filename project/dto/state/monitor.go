package state

type Monitor struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Workspace Workspace `json:"activeWorkspace"`
	Scale     float64   `json:"scale"`
	Focused   bool      `json:"focused"`
}

package state

type Window struct {
	Address   string    `json:"address"`
	Class     string    `json:"class"`
	Floating  bool      `json:"floating"`
	Location  Vec2      `json:"at"`
	MonitorId int       `json:"monitor"`
	Pid       int       `json:"pid"`
	Size      Vec2      `json:"size"`
	Tags      []string  `json:"tags"`
	Title     string    `json:"title"`
	Workspace Workspace `json:"workspace"`
}

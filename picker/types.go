package picker

type EventCaptcha struct {
	DropID string    `json:"drop_id"`
	Codes  [2]string `json:"codes"`
}

type EventPicked struct {
	DropID string `json:"drop_id,omitempty"`
}

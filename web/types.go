package web

type channel struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type user struct {
	ID     string `json:"id,omitempty"`
	Nick   string `json:"nick,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type pick struct {
	ID     string `json:"id,omitempty"`
	DropID string `json:"drop_id,omitempty"`
	Code   string `json:"code,omitempty"`
	Since  string `json:"since,omitempty"`
	User   *user  `json:"user,omitempty"`
}

type captcha struct {
	DropID string   `json:"drop_id,omitempty"`
	Codes  []string `json:"codes,omitempty"`
}

type drop struct {
	ID        string   `json:"id,omitempty"`
	Currency  int      `json:"currency,omitempty"`
	Timestamp int64    `json:"timestamp,omitempty"`
	Captcha   *captcha `json:"captcha,omitempty"`
	Active    *bool    `json:"active,omitempty"`
	Since     *string  `json:"since,omitempty"`
	Winner    *user    `json:"winner,omitempty"`
	Picks     []pick   `json:"picks,omitempty"`
}

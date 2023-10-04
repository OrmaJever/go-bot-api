package models

type Bot struct {
	Id       int64    `json:"id"`
	Token    string   `json:"token"`
	Secret   string   `json:"secret"`
	AdminId  int64    `json:"admin_id"`
	Packages []string `json:"packages"`
}

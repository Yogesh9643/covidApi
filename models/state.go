package models

type StateList struct {
	StateList []State `json:"statewise"`
}

type State struct {
	Active          string `json:"active"`
	Confirmed       string `json:"confirmed"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	State           string `json:"state"`
}

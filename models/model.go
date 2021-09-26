package models

//Coordinate is a [longitude, latitude]
type Statelist struct {
	Statelist []State `json:"statewise"`
}

type State struct {
	Active          string `json:"active"`
	Confirmed       string `json:"confirmed"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	State           string `json:"state"`
}

type Responsejson struct {
	Active          string `json:"active"`
	Confirmed       string `json:"confirmed"`
	Lastupdatedtime string `json:"lastupdatedtime"`
	State           string `json:"state"`
	Totalincountry  string `json:"totalconfirmedindia"`
}

type StateMapBox struct {
	StateMapBox []StateCoor `json:"features"`
}

type StateCoor struct {
	StateName string `json:"text"`
}

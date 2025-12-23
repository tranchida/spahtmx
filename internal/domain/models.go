package domain

type User struct {
	ID       string
	Username string
	Email    string
	Status   bool
}

type PrizeList struct {
	Prizes []Prize `json:"prizes"`
}

type Prize struct {
	ID                string     `json:"id"`
	Year              string     `json:"year"`
	Category          string     `json:"category"`
	OverallMotivation string     `json:"overallMotivation,omitempty"`
	Laureates         []Laureate `json:"laureates,omitempty"`
}

type Laureate struct {
	Firstname  string `json:"firstname,omitempty"`
	Surname    string `json:"surname,omitempty"`
	Motivation string `json:"motivation,omitempty"`
	Share      string `json:"share,omitempty"`
}

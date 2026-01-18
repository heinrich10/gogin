package model

type Person struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CountryCode string `json:"country_code"`
	UpdatedAt   string `json:"update_at"`  // timestamp of the last update
	CreatedAt   string `json:"created_at"` // timestamp of creation
}

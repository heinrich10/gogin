package model

type Person struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name" binding:"required,min=1,max=255"`
	LastName    string `json:"last_name" binding:"max=255"`
	CountryCode string `json:"country_code" binding:"required,len=2"`
	UpdatedAt   string `json:"update_at,omitzero"`  // timestamp of the last update
	CreatedAt   string `json:"created_at,omitzero"` // timestamp of creation
}

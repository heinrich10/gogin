package model

type Country struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Phone         int64  `json:"phone"`
	Symbol        string `json:"symbol"`
	Capital       string `json:"capital"`
	Currency      string `json:"currency"`
	ContinentCode string `json:"continent_code"`
	Alpha3        string `json:"alpha3"`
	updateAt      int64  // timestamp of the last update
	createdAt     int64  // timestamp of creation
}

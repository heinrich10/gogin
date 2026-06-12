package model

type Country struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Symbol        string `json:"symbol"`
	Capital       string `json:"capital"`
	Currency      string `json:"currency"`
	ContinentCode string `json:"continent_code"`
	Alpha3        string `json:"alpha_3"`
}

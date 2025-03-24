package dtos

type CreateBloodPressure struct {
	Systiolic int    `json:"systiolic"`
	Diastolic int    `json:"diastolic"`
	Pulse     int    `json:"pulse"`
	Medicine  string `json:"medicine"`
}
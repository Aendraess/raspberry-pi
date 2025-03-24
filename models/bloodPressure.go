package models

type BloodPressure struct {
	BaseModel
	Systiolic int    `json:"systiolic"`
	Diastolic int    `json:"diastolic"`
	Pulse     int    `json:"pulse"`
	Medicine  string `json:"medicine"`
}
package dtos

type CreateBloodPressure struct {
	Systolic  int    `json:"systolic"`
	Diastolic int    `json:"diastolic"`
	Pulse     int    `json:"pulse"`
	Medicine  string `json:"medicine"`
}
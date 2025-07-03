package pageGeneration

import "time"

type ControlQuadrantData struct {
	QuadrantData
	CurrentDate string
}

func (pg *PageGenerator) GenerateControlQuadrant() (*ControlQuadrantData, error) {
	currentDate := time.Now().Format("2006-01-02")

	return &ControlQuadrantData{
		QuadrantData: QuadrantData{Title: "Control"},
		CurrentDate:  currentDate,
	}, nil
}

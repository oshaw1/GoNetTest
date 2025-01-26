package pageGeneration

type ControlQuadrantData struct {
	QuadrantData
}

func (pg *PageGenerator) GenerateControlQuadrant() (*ControlQuadrantData, error) {
	return &ControlQuadrantData{
		QuadrantData: QuadrantData{Title: "Control"},
	}, nil
}

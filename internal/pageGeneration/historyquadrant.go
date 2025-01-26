package pageGeneration

type GenerateQuadrantData struct {
	QuadrantData
}

func (pg *PageGenerator) GenerateHistoryQuadrant() (*GenerateQuadrantData, error) {
	return &GenerateQuadrantData{
		QuadrantData: QuadrantData{Title: "History"},
	}, nil
}

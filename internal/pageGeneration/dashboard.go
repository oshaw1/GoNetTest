package pageGeneration

import "net/http"

func (pg *PageGenerator) RenderDashboard(w http.ResponseWriter) error {
	testData, err := pg.GenerateTestQuadrant("", "")
	if err != nil {
		return err
	}
	generateData, err := pg.GenerateHistoryQuadrant()
	if err != nil {
		return err
	}
	controlData, err := pg.GenerateControlQuadrant()
	if err != nil {
		return err
	}
	schedulerData, err := pg.GenerateSchedulerQuadrant()
	if err != nil {
		return err
	}
	data := &DashboardData{
		TestData:      testData,
		GenerateData:  generateData,
		ControlData:   controlData,
		SchedulerData: schedulerData,
	}
	return pg.templates.ExecuteTemplate(w, "base", data)
}

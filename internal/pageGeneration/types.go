package pageGeneration

type QuadrantData struct {
	Title string
	Error error
}

type DashboardData struct {
	TestData      *TestQuadrantData
	ControlData   *ControlQuadrantData
	SchedulerData *SchedulerQuadrantData
}

type TestGroup struct {
	TimeGroup  string
	JsonPath   string
	TestResult interface{}
	ChartPaths map[string]string
}

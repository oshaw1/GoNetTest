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
	ResultID   int64
	ChartIDs   string // comma-separated charts.id list; only set for Historic groups
	TestResult interface{}
	ChartPaths map[string]string
	Historic   bool
}

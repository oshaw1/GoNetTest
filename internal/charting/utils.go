package charting

import "github.com/go-echarts/go-echarts/v2/opts"

func generateBarItems(speeds []float64) []opts.BarData {
	items := make([]opts.BarData, len(speeds))
	for i, speed := range speeds {
		items[i] = opts.BarData{Value: speed}
	}
	return items
}

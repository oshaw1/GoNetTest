package charting

import "github.com/go-echarts/go-echarts/v2/opts"

func generateBarItems(speeds []float64) []opts.BarData {
	items := make([]opts.BarData, len(speeds))
	for i, speed := range speeds {
		items[i] = opts.BarData{Value: speed}
	}
	return items
}

func generateLineItems(data []float64) []opts.LineData {
	items := make([]opts.LineData, 0, len(data))
	for _, d := range data {
		items = append(items, opts.LineData{Value: d})
	}
	return items
}

func calculateAverage(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func findMin(values []float64) float64 {
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func findMax(values []float64) float64 {
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

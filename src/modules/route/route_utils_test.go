package route_test

import (
	"sort"
	"testing"
)

var (
	distanceMatrix = [][]float64{
		{0.00, 42.60, 47.50, 46.40, 53.00, 49.20, 63.20, 51.60, 46.40, 47.70, 36.80, 39.70},
		{42.60, 0.00, 9.80, 3.50, 19.70, 15.20, 14.90, 11.20, 4.90, 13.60, 7.30, 6.50},
		{47.50, 9.80, 0.00, 8.20, 13.50, 5.80, 26.50, 15.50, 7.80, 2.70, 15.20, 7.60},
		{46.40, 3.50, 8.20, 0.00, 24.70, 14.40, 14.10, 10.40, 2.70, 11.90, 9.70, 9.00},
		{53.00, 19.70, 13.50, 24.70, 0.00, 7.60, 41.90, 36.40, 20.70, 13.00, 25.90, 18.90},
		{49.20, 15.20, 5.80, 14.40, 7.60, 0.00, 32.10, 21.10, 13.40, 6.00, 18.30, 10.70},
		{63.20, 14.90, 26.50, 14.10, 41.90, 32.10, 0.00, 5.30, 12.70, 21.60, 19.50, 19.50},
		{51.60, 11.20, 15.50, 10.40, 36.40, 21.10, 5.30, 0.00, 8.50, 18.50, 14.00, 16.30},
		{46.40, 4.90, 7.80, 2.70, 20.70, 13.40, 12.70, 8.50, 0.00, 10.50, 19.70, 13.10},
		{47.70, 13.60, 2.70, 11.90, 13.00, 6.00, 21.60, 18.50, 10.50, 0.00, 15.50, 8.00},
		{36.80, 7.30, 15.20, 9.70, 25.90, 18.30, 19.50, 14.00, 19.70, 15.50, 0.00, 9.90},
		{39.70, 6.50, 7.60, 9.00, 18.90, 10.70, 19.50, 16.30, 13.10, 8.00, 9.90, 0.00},
	}
	pointCapacities = []float64{700.00, 700.00, 586.00, 714.00, 729.00, 529.00, 571.00, 529.00, 571.00, 629.00, 714.00}
	vehicleCapacity = []float64{4000.00, 4000.00, 4000.00}
)

func testClarkeWrightSaving(distanceMatrix [][]float64, pointCapacities []float64, vehicleCapacity []float64) [][]int {
	numPoints := len(distanceMatrix)

	savings := [][]interface{}{}
	for i := 1; i < numPoints; i++ {
		for j := 1; j < i; j++ {
			saving := distanceMatrix[0][i] + distanceMatrix[0][j] - distanceMatrix[i][j]
			savings = append(savings, []interface{}{i - 1, j - 1, saving})
		}
	}

	sort.Slice(savings, func(i, j int) bool {
		return savings[i][2].(float64) > savings[j][2].(float64)
	})

	routes := [][]int{}
	carriedCapacity := 0.0
	vehicleIndex := 0
	visitedPoint := []int{}
	for idx := range savings {
		firstPointIndex := testCalculatePointVisited(&visitedPoint, savings[idx][0].(int))
		secondPointIndex := testCalculatePointVisited(&visitedPoint, savings[idx][1].(int))
		if secondPointIndex >= 0 {
			continue
		}
		if firstPointIndex >= 0 && firstPointIndex != len(visitedPoint)-1 {
			continue
		}

		if !(idx > 0 && savings[idx][0].(int) == savings[idx-1][1].(int)) {
			if carriedCapacity+pointCapacities[savings[idx][0].(int)]+pointCapacities[savings[idx][1].(int)] > vehicleCapacity[vehicleIndex] {
				vehicleIndex++
				carriedCapacity = 0
			}
		} else {
			if carriedCapacity+pointCapacities[savings[idx][0].(int)] > vehicleCapacity[vehicleIndex] {
				vehicleIndex++
				carriedCapacity = 0
			}
		}

		if len(routes)-1 < vehicleIndex {
			routes = append(routes, []int{})
		}

		if !(idx > 0 && savings[idx][0].(int) == savings[idx-1][1].(int)) {
			visitedPoint = append(visitedPoint, savings[idx][0].(int))
			routes[vehicleIndex] = append(routes[vehicleIndex], savings[idx][0].(int))
			carriedCapacity += pointCapacities[savings[idx][0].(int)]
		}

		visitedPoint = append(visitedPoint, savings[idx][1].(int))
		routes[vehicleIndex] = append(routes[vehicleIndex], savings[idx][1].(int))
		carriedCapacity += pointCapacities[savings[idx][1].(int)]
	}

	return routes
}

func testCalculatePointVisited(visitedPoint *[]int, point int) int {
	index := -1
	for i := range *visitedPoint {
		if (*visitedPoint)[i] == point {
			index = i
			break
		}
	}
	return index
}

func TestClarkeWrightSaving(t *testing.T) {
	correctRoutes := [][]int{{6, 5, 2, 4, 3}, {8, 1, 7, 0, 10, 9}}
	generatedRoutes := testClarkeWrightSaving(distanceMatrix, pointCapacities, vehicleCapacity)

	if len(correctRoutes) != len(generatedRoutes) {
		t.Fatalf("len(correctRoutes): %d\nlen(generatedRoutes): %d", len(correctRoutes), len(generatedRoutes))
	}

	isCorrect := true

	for i := range generatedRoutes {
		for j := range generatedRoutes[i] {
			if generatedRoutes[i][j] != correctRoutes[i][j] {
				isCorrect = false
				break
			}
		}
		if !isCorrect {
			break
		}
	}

	if !isCorrect {
		t.Fatalf("correctRoutes: %v\ngeneratedRoutes: %v", correctRoutes, generatedRoutes)
	}
}

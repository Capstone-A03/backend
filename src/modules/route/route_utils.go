package route

import (
	"math"
	"sort"
)

type RouteNode struct {
	Coordinate *RouteCoordinate
	Capacity   *float64
}

type RouteCoordinate struct {
	Latitude  *float64
	Longitude *float64
}

func clarkeWrightSaving(routeNodes *[]*RouteNode, vehiclesCapacity *[]float64) *[][]int {
	numPoints := len(*routeNodes) + 1
	distanceMatrix := make([][]float64, numPoints)
	for i := 0; i < numPoints; i++ {
		distanceMatrix[i] = make([]float64, numPoints)
	}
	for i := 0; i < numPoints; i++ {
		for j := 0; j < i; j++ {
			distanceMatrix[i][j] = calculateDistance((*routeNodes)[j].Coordinate, (*routeNodes)[i].Coordinate)
		}
	}

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

	route := [][]int{}
	carriedCapacity := 0.0
	vehicleIndex := 0
	visitedPoint := []int{}
	for idx := range savings {
		firstPointIndex := calculatePointVisited(&visitedPoint, savings[idx][0].(int))
		secondPointIndex := calculatePointVisited(&visitedPoint, savings[idx][1].(int))
		if secondPointIndex >= 0 {
			continue
		}
		if firstPointIndex >= 0 && firstPointIndex != len(visitedPoint)-1 {
			continue
		}

		if !(idx > 0 && savings[idx][0].(int) == savings[idx-1][1].(int)) {
			if carriedCapacity+*((*routeNodes)[savings[idx][0].(int)].Capacity)+*((*routeNodes)[savings[idx][1].(int)].Capacity) > (*vehiclesCapacity)[vehicleIndex] {
				vehicleIndex++
				carriedCapacity = 0
			}
		} else {
			if carriedCapacity+*((*routeNodes)[savings[idx][0].(int)].Capacity) > (*vehiclesCapacity)[vehicleIndex] {
				vehicleIndex++
				carriedCapacity = 0
			}
		}

		if len(route)-1 < vehicleIndex {
			route = append(route, []int{})
		}

		if !(idx > 0 && savings[idx][0].(int) == savings[idx-1][1].(int)) {
			visitedPoint = append(visitedPoint, savings[idx][0].(int))
			route[vehicleIndex] = append(route[vehicleIndex], savings[idx][0].(int))
			carriedCapacity += *((*routeNodes)[savings[idx][0].(int)].Capacity)
		}

		visitedPoint = append(visitedPoint, savings[idx][1].(int))
		route[vehicleIndex] = append(route[vehicleIndex], savings[idx][1].(int))
		carriedCapacity += *((*routeNodes)[savings[idx][1].(int)].Capacity)
	}

	return &route
}

func calculateDistance(coordinate1, coordinate2 *RouteCoordinate) float64 {
	return math.Sqrt(math.Pow(*(coordinate2.Latitude)-*(coordinate1.Latitude), 2) + math.Pow(*(coordinate2.Longitude)-*(coordinate1.Longitude), 2))
}

func calculatePointVisited(visitedPoint *[]int, point int) int {
	index := -1
	for i := range *visitedPoint {
		if (*visitedPoint)[i] == point {
			index = i
			break
		}
	}
	return index
}

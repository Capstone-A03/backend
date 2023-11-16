package route

import (
	"math"
	"sort"

	"github.com/google/uuid"
)

type VehicleCapacity struct {
	ID             *uuid.UUID
	Capacity       *float64
	FilledCapacity *float64
}

type RouteNode struct {
	ID                     *uuid.UUID
	Coordinate             *RouteCoordinate
	RemainingWasteCapacity *float64
}

type RouteCoordinate struct {
	Latitude  *float64
	Longitude *float64
}

type generatedRoute struct {
	TruckID *uuid.UUID
	DumpIDs *[]*uuid.UUID
}

func clarkeWrightSaving(routeNodes *[]*RouteNode, vehiclesCapacity *[]*VehicleCapacity) *[]*generatedRoute {
	sort.Slice(*vehiclesCapacity, func(i, j int) bool {
		return *(*vehiclesCapacity)[i].Capacity > *(*vehiclesCapacity)[j].Capacity
	})

	numPoints := len(*routeNodes)
	distanceMatrix := make([][]float64, numPoints)
	for i := 0; i < numPoints; i++ {
		distanceMatrix[i] = make([]float64, numPoints)
		for j := 0; j < i; j++ {
			distanceMatrix[i][j] = calculateDistance((*routeNodes)[j].Coordinate, (*routeNodes)[i].Coordinate)
		}
	}

	savings := [][]interface{}{}
	for i := 1; i < numPoints; i++ {
		for j := 1; j < i; j++ {
			saving := distanceMatrix[0][i] + distanceMatrix[0][j] - distanceMatrix[i][j]
			savings = append(savings, []interface{}{i, j, saving})
		}
	}

	sort.Slice(savings, func(i, j int) bool {
		return savings[i][2].(float64) > savings[j][2].(float64)
	})

	route := new([]*generatedRoute)
	vehicleIndex := 0
	visitedPoint := []int{}

savingsLoop:
	for idx := range savings {
		firstPointIndex := findPointVisitedIndex(&visitedPoint, savings[idx][0].(int))
		secondPointIndex := findPointVisitedIndex(&visitedPoint, savings[idx][1].(int))
		if secondPointIndex >= 0 {
			continue
		}
		if firstPointIndex >= 0 && firstPointIndex != len(visitedPoint)-1 {
			continue
		}

		for *((*routeNodes)[savings[idx][0].(int)].RemainingWasteCapacity) > 0 {
			if vehicleIndex > len(*vehiclesCapacity)-1 {
				break savingsLoop
			}

			garbageToBeCollectedVolume := *((*routeNodes)[savings[idx][0].(int)].RemainingWasteCapacity)
			truckRemainingCapacity := *(*vehiclesCapacity)[vehicleIndex].Capacity - *(*vehiclesCapacity)[vehicleIndex].FilledCapacity
			if garbageToBeCollectedVolume > truckRemainingCapacity {
				garbageToBeCollectedVolume = truckRemainingCapacity
			}
			*(*vehiclesCapacity)[vehicleIndex].FilledCapacity += garbageToBeCollectedVolume
			*(*routeNodes)[savings[idx][0].(int)].RemainingWasteCapacity -= garbageToBeCollectedVolume

			if *(*vehiclesCapacity)[vehicleIndex].FilledCapacity >= *(*vehiclesCapacity)[vehicleIndex].Capacity*85/100 {
				vehicleIndex++
			}

			for len(*route)-1 < vehicleIndex {
				*route = append(*route, &generatedRoute{
					TruckID: (*vehiclesCapacity)[vehicleIndex].ID,
					DumpIDs: &[]*uuid.UUID{},
				})
			}

			if !(idx > 0 && savings[idx][0].(int) == savings[idx-1][1].(int)) {
				visitedPoint = append(visitedPoint, savings[idx][0].(int))
				if len(*(*route)[vehicleIndex].DumpIDs) == 0 {
					*(*route)[vehicleIndex].DumpIDs = append(*(*route)[vehicleIndex].DumpIDs, (*routeNodes)[0].ID)
				}
				*(*route)[vehicleIndex].DumpIDs = append(*(*route)[vehicleIndex].DumpIDs, (*routeNodes)[savings[idx][0].(int)].ID)
			}

			visitedPoint = append(visitedPoint, savings[idx][1].(int))
			if len(*(*route)[vehicleIndex].DumpIDs) == 0 {
				*(*route)[vehicleIndex].DumpIDs = append(*(*route)[vehicleIndex].DumpIDs, (*routeNodes)[0].ID)
			}
			*(*route)[vehicleIndex].DumpIDs = append(*(*route)[vehicleIndex].DumpIDs, (*routeNodes)[savings[idx][1].(int)].ID)
		}
	}

	return route
}

func calculateDistance(coordinate1, coordinate2 *RouteCoordinate) float64 {
	return math.Sqrt(math.Pow(*(coordinate2.Latitude)-*(coordinate1.Latitude), 2) + math.Pow(*(coordinate2.Longitude)-*(coordinate1.Longitude), 2))
}

func findPointVisitedIndex(visitedPoint *[]int, point int) int {
	index := -1
	for i := range *visitedPoint {
		if (*visitedPoint)[i] == point {
			index = i
			break
		}
	}
	return index
}

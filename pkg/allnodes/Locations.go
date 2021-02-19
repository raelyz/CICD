package allnodes

import (
	"fmt"
	"sync"
)

//Locations saves locations as a slice with the currentSize and maxSize known. This will be inherited by most of the other structs such as Itinerary, LocationRequest and DatabaseLocations
type Locations struct {
	locations   []*Location
	currentSize int
	Mutex       sync.RWMutex
}

//AddLocation adds a new location to the last last empty slot in the slice.
func (l *Locations) AddLocation(newLocation **Location) {
	if l.currentSize == 10 {
		fmt.Println("full capacity")
	} else {
		switch l.currentSize {
		case 0:
			l.locations = make([]*Location, 10)
			l.locations[0] = *newLocation
		default:
			l.locations[l.currentSize] = *newLocation
		}
		(*newLocation).Popularity++
		l.currentSize++
	}

}

//RemoveLocation removes the location from its specific index and returns the new slice without the element.
func (l *Locations) RemoveLocation(index int) error {
	if index > 0 && index <= l.currentSize {
		fmt.Println("Removing", l.locations[index-1].name, "from index", index)
		if index == l.currentSize {
			l.locations[index-1] = &Location{}
		} else if index < l.currentSize {
			l.locations[index-1].DecreasePopularity()
			for i := index - 1; i < l.currentSize-1; i++ {
				l.locations[i] = l.locations[i+1]
			}
		}
		l.currentSize--
		return nil
	}
	return fmt.Errorf("Invalid Index")

}

//PrintLocations prints l Locations in the order they were added
func (l *Locations) PrintLocations() {
	if l.currentSize == 0 {
		fmt.Println("No Locations loaded in the database, please add one before viewing")
		return
	}
	fmt.Println("==================================================")
	for i := 0; i < l.currentSize; i++ {
		fmt.Println(i+1, (*l.locations[i]).name, (*l.locations[i]).Popularity)
	}
}

func (l *Locations) merge(firstHalf []*Location, secondHalf []*Location) []*Location {
	var i, j, index, arrayLength int
	arrayLength = len(firstHalf) + len(secondHalf)
	sortedArray := make([]*Location, arrayLength)
	for i < len(firstHalf) || j < len(secondHalf) {
		if i >= len(firstHalf) {
			sortedArray[index] = secondHalf[j]
			j++
			index++
		} else if j >= len(secondHalf) {
			sortedArray[index] = firstHalf[i]
			i++
			index++
		} else {
			if (*firstHalf[i]).Popularity < (*secondHalf[j]).Popularity {
				sortedArray[index] = secondHalf[j]
				j++
			} else {
				sortedArray[index] = firstHalf[i]
				i++
			}
			index++
		}
	}
	return sortedArray
}

func (l *Locations) mergeSort(data []*Location, size int) []*Location {
	if size < 2 {
		return data
	}
	mid := size / 2

	firstHalf := make([]*Location, mid)
	for i := 0; i < mid; i++ {
		firstHalf[i] = data[i]
	}
	firstHalf = l.mergeSort(firstHalf, mid)
	secondHalf := make([]*Location, size-mid)
	var j int = mid
	for i := 0; i < size-mid; i++ {
		secondHalf[i] = data[j]
		j++
	}
	secondHalf = l.mergeSort(secondHalf, size-mid)

	return l.merge(firstHalf, secondHalf)
}

//SortByPopularity sorts all locations based ober of times someone has added it into their itinerary
func (l *Locations) SortByPopularity() {
	l.locations = l.mergeSort(l.locations, l.currentSize)
	l.PrintLocations()
}

//Contains prints l Locations in the order they were added
func (l *Locations) Contains(locat **Location) bool {
	for i := 0; i < l.currentSize; i++ {
		if l.locations[i] == *locat {
			return true
		}
	}
	return false
}

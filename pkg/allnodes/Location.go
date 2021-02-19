package allnodes

import "sync"

// Location represents one particular sightseeing destination.
type Location struct {
	name        string
	Description string
	Popularity  int
	Mutex       sync.Mutex
}

//CreateLocation given the name and returns it
func CreateLocation(name string, description string, rating int) *Location {
	var mutex sync.Mutex
	return &Location{name, description, rating, mutex}
}

// PrintInfo prints the location info and its rating
func (l *Location) String() string {
	return l.name
}

//IncreasePopularity is called when a user has added the Location into his Itinerary
func (l *Location) IncreasePopularity() {
	l.Mutex.Lock()
	l.Popularity++
	l.Mutex.Unlock()
}

//DecreasePopularity is called when it has been removed from an Itinerary
func (l *Location) DecreasePopularity() {
	l.Mutex.Lock()
	l.Popularity--
	l.Mutex.Unlock()
}

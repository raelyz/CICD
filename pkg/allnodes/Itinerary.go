package allnodes

//Itinerary saves all the locations the user has added into an array of locations and the number of locations are noted. The route saved is the optimized proposed route with the traveltime indicated.
type Itinerary struct {
	route      *ItineraryNode
	travelTime float64
	Locations
}

//ItineraryNode refers to each Location that is stored within the Itinerary. Visited is tracked to help the users gauge if they are still on schedule.
type ItineraryNode struct {
	curr    *Location
	next    *ItineraryNode
	visited bool
}

//NewItinerary returns an empty Itinerary address
func NewItinerary() *Itinerary {
	return &Itinerary{}
}

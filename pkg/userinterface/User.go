package userinterface

import (
	"Assignment2/pkg/allnodes"
	"errors"
	"fmt"
	"strconv"
)

//UserMenu contains the information for User's View
type UserMenu struct {
	BasicMenu
	itinerary *allnodes.Itinerary
}

//NewUserMenu returns the address of an empty menu for utilization
func NewUserMenu() *UserMenu {
	return &UserMenu{}
}

//SeedData Populates the Menu with the Options
func (u *UserMenu) SeedData() {
	u.options = []string{"Browse Locations", "Print Itinerary", "Add to Itinerary", "Remove from Itinerary", "sort itinerary by popularity", "Exit"}
	u.itinerary = allnodes.NewItinerary()
	// l.Users.SeedData()
}

//MenuOptions receives the user selected input and runs the respective method calls
func (u *UserMenu) MenuOptions() int {
	for {
		u.printMenu()
		input := u.userSelection()
		switch input {
		case 1:
			u.printAllLocations()
		case 2:
			u.printItinerary()
		case 3:
			u.addNewLocation()
		case 4:
			u.deleteLocation()
		case 5:
			u.itinerary.SortByPopularity()
		case 6:
			fmt.Println("Quitting Application See you Again Next Time")
			return 0
		}
	}
}

func (u *UserMenu) printItinerary() {
	u.itinerary.PrintLocations()
}

//prints all locations stored in the database
func (u *UserMenu) printAllLocations() {
	mutex.Lock()
	firstNode := (data.FirstNode())
	var j int = 1
	fmt.Println("==================================================")
	for firstNode != nil {
		for i := 0; i < (*firstNode).NumKeys; i++ {
			fmt.Println(j, (*firstNode).Pointers[i], ",Popularity:", (*firstNode).Pointers[i].(*allnodes.Location).Popularity, ",UID:", (*firstNode).Keys[i], ",Description:", (*firstNode).Pointers[i].(*allnodes.Location).Description)
			j++
		}
		firstNode = (*firstNode).Next
	}
	fmt.Println("==================================================")
	mutex.Unlock()
}

func (u *UserMenu) addNewLocation() {
	var input string
	errMessage := "Please key in a valid ID"
	fmt.Println("Please key in the Location ID to be Added")
	fmt.Println("==================================================")
	fmt.Scanln(&input)
	if index, err := strconv.Atoi(input); err != nil {
		fmt.Println(errors.New(errMessage))
	} else {
		location, err := data.Search(index)
		if err != nil {
			fmt.Println(err)
		} else {
			if u.itinerary.Contains(&location) {
				fmt.Printf("Location %s exists in the Itinerary \n", location.String())
			} else {
				u.itinerary.AddLocation(&location)
			}
		}
	}

}

func (u *UserMenu) deleteLocation() {
	var input string
	errMessage := "Please key in a ID"
	fmt.Println("Please key in the ID to be deleted")
	fmt.Scanln(&input)
	if index, err := strconv.Atoi(input); err != nil {
		fmt.Println(errors.New(errMessage))
	} else {
		err := u.itinerary.RemoveLocation(index)
		if err != nil {
			fmt.Println(err)
		}
	}
}

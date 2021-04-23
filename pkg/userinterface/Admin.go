package userinterface

import (
	"pkg/allnodes"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
)

var (
	data             = allnodes.NewBPlusTree()
	savedItineraries = make([]allnodes.Itinerary, 10)
	mutex            sync.Mutex
)

//LoadData initializes the data for users to utilize
func LoadData() {
	data.SeedData()
	loadItineraries()
}

//loadItineraries loads Itineraries that other people have planned
func loadItineraries() {

}

//AdminMenu contains the information for Administrator View
type AdminMenu struct {
	BasicMenu
	mutex sync.Mutex
}

//NewAdminMenu returns the address of an empty menu for utilization
func NewAdminMenu() *AdminMenu {
	return &AdminMenu{}
}

//SeedData Populates the Menu with the Options
func (a *AdminMenu) SeedData() {
	a.options = []string{"View All Locations", "Add New Location", "Delete Location", "Exit"}
}

//MenuOptions receives the user selected input and runs the respective method calls
func (a *AdminMenu) MenuOptions() int {
	for {
		a.printMenu()
		input := a.userSelection()
		switch input {
		case 1:
			a.printAllLocations()
		case 2:
			a.addNewLocation()
		case 3:
			a.deleteLocation()
		case 4:
			fmt.Println("Quitting Application See you Again Next Time")
			return 0
		}
	}
}

//prints all locations stored in the database
func (a *AdminMenu) printAllLocations() {
	firstNode := (data.FirstNode())
	fmt.Println("==================================================")
	for firstNode != nil {
		for i := 0; i < (*firstNode).NumKeys; i++ {
			fmt.Println("Name:", (*firstNode).Pointers[i], ",Popularity:", (*firstNode).Pointers[i].(*allnodes.Location).Popularity, ",UID:", (*firstNode).Keys[i], ",Description:", (*firstNode).Pointers[i].(*allnodes.Location).Description)

		}
		firstNode = (*firstNode).Next
	}
	fmt.Println("==================================================")
}

func (a *AdminMenu) addNewLocation() {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("What is the venue called?")
	fmt.Println("==================================================")
	text, _, _ := inputReader.ReadLine()
	location := string(text)
	fmt.Println("Please provide a description of the location")
	fmt.Println("==================================================")
	text, _, _ = inputReader.ReadLine()
	description := string(text)
	mutex.Lock()
	data.InsertTree(allnodes.CreateLocation(location, description, 0))
	mutex.Unlock()
}

func (a *AdminMenu) deleteLocation() {
	var input string
	errMessage := "Please key in a valid Number"
	fmt.Println("Please key in the ID to be deleted")
	fmt.Println("==================================================")
	fmt.Scanln(&input)
	if index, err := strconv.Atoi(input); err != nil {
		fmt.Println(errors.New(errMessage))
		fmt.Println("==================================================")
	} else {
		mutex.Lock()
		_, err := data.Delete(index)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Item index %v deleted successfully \n", index)
		}
		mutex.Unlock()
	}
}

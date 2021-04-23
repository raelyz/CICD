package main

import (
	"pkg/allnodes"
	"pkg/userinterface"
	"sync"
)

var (
	itinerary = allnodes.NewItinerary()
	login     = userinterface.NewLoginMenu()
	admin     = userinterface.NewAdminMenu()
	user      = userinterface.NewUserMenu()
	wg        sync.WaitGroup
)

func main() {
	userinterface.LoadData()
	login.SeedData()
	admin.SeedData()
	user.SeedData()
	input := login.MenuOptions()

	// Possible Location to run goroutines, would require passing the address of the waitgroup through but didn't have time to test for concurrency stability
	// wg.Add(2)
	// go user.MenuOptions()
	// go user.MenuOptions()
	// wg.Wait()
	//
	if input == 1 {
		user.MenuOptions()
	} else if input == 2 {
		admin.MenuOptions()
	}

}

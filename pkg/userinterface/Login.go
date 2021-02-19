package userinterface

import (
	"fmt"
)

//LoginMenu contains the methods pertaining to Logging in / Account Creation and UserData
type LoginMenu struct {
	BasicMenu
	failedLogin int
}

//NewLoginMenu returns the address of an empty menu for utilization
func NewLoginMenu() *LoginMenu {
	return &LoginMenu{}
}

//SeedData Populates the Menu with the Options
func (l *LoginMenu) SeedData() {
	l.options = []string{"User", "Admin", "Exit"}
	// l.Users.SeedData()
}

//MenuOptions receives the user selected input and runs the respective method calls
func (l *LoginMenu) MenuOptions() int {
	for {
		l.printMenu()
		input := l.userSelection()
		if l.failedLogin == 3 {
			input = 3
			fmt.Println("Too many attempts, bye")
		}
		switch input {
		case 1:
			return 1
		case 2:
			ok := l.requestPassword()
			if ok {
				return 2
			}
		case 3:
			fmt.Println("Quitting Application See you Again Next Time")
			return 0
		}
	}
}

func (l *LoginMenu) requestPassword() bool {
	var input string
	fmt.Println("Please key in the administator password")
	fmt.Scanln(&input)
	if input == "123" {
		return true
	}
	l.failedLogin++
	fmt.Println("Failed login attempt.", 3-l.failedLogin, " attempts left")
	return false
}

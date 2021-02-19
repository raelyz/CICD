package userinterface

import (
	"Assignment2/pkg/allnodes"
	"fmt"
	"strconv"
)

//BasicMenu UI that stores the options as a string and it includes basic functionality like input validation and printing options
type BasicMenu struct {
	options   []string
	locations *allnodes.BPlusTree
	Menu
}

//Menu Interface for modularity
type Menu interface {
	PrintMenu()
	MenuInput()
	UserSelection()
}

//PrintMenu prints all the options for the given MenuType
func (b *BasicMenu) printMenu() {

	fmt.Println("Please Select from one of the following")
	fmt.Println("==================================================")
	for index, options := range b.options {
		fmt.Println(index+1, options)
	}
}

//MenuInput validates that the input is indeed an integer and is within the range of the available choices
func (b *BasicMenu) menuInput(input string) (int, error) {
	if input, err := strconv.Atoi(input); err != nil {
		return 0, err
	} else {
		if input > 0 && input <= len(b.options) {
			return input, nil
		}

	}
	return 0, nil
}

//UserSelection prompts the user continuiously until a valid choice is made and returns that value
func (b *BasicMenu) userSelection() int {
	for {
		var input string
		fmt.Scanln(&input)
		if output, err := b.menuInput(input); err != nil {
			fmt.Println("Please key in a valid input")
			continue
		} else {
			return output
		}
	}
}

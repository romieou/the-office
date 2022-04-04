package main

import (
	"os"
	"bufio"
	"log"
	"fmt"
	"strings"
	"v0/office"
	"errors"
)

func main() {
	// open employeeActions and return error, if any
	file, err := os.Open("employeeActions")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// parse lines and add them to a 2D slice
	scanner := bufio.NewScanner(file)
	var lines [][]string
	for scanner.Scan() {
		lines = append(lines, strings.Split(strings.ReplaceAll(scanner.Text(), " ", ""), ","))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	// keeps track of employee's location
	var currRoom office.Location
	// declare a list of employees which will store all employee data
	var listOfEmployees []*office.Employee
	// set employee count
	i := 0
	// iterate through lines and add employees to listOfEmployees
	for j, line := range lines {
		// checker to check if everything went right
		ok := true
		// if line is not empty, take first string and convert it into an employee of respective title
		if len(line) >= 1 {
		// print error if failed to create employee and move on to next line
			newEmployee, err := office.NewEmployeeFactory(line[0]) 
			if err != nil {
				title := ""
				if errors.Is(err, office.ErrUnknownEmplType) {
					title = line[0]
				}
				fmt.Printf("\nemployeeActions Line %d: fail: <main.creating newEmployee>: %s %q\n\n", j, err.Error(), title)
				continue
			}
			// append each employee to employee list
			listOfEmployees = append(listOfEmployees, &newEmployee)
		// if there's no error, increment employee count by one (that's how we will be referring to each employee in standard output)
			i++
			fmt.Printf("********************Employee %d *********************\n", i)
			fmt.Println("title: ", office.GetEmployeeTitle(newEmployee), "\n")
		// remove employee title from line and leave plain locations
			line = line[1:]
		// iterate through line (which now stores locations only) and create new location for each room 
			for _, l := range line {
				newLocation, err := office.NewLocationFactory(l)
				if err != nil {
					title := ""
					if errors.Is(err, office.ErrUnknownLocationType) {
						title = l
					}
					fmt.Printf("\nemployeeActions Line %d: fail: <main.creating newLocation>: %s %q\n\n", j, err.Error(), title)
		// error has occured, so cannot print out success message
					ok = false
					break
				}
		// if there's no error, print the current location of the employee
				currLocation:= newEmployee.GetCurrentLocation()
				fmt.Println("I'm at the ... ", office.GetLocationName(currLocation))
		// move the employee to the next room; if not possible, break and move to the next employee 
				if err := newEmployee.MoveToLocation(newLocation); err != nil {
					fmt.Printf("\nemployeeActions Line %d: fail: <main.Movetolocation>: %s\n\n", j, err)
					ok = false
					break
				}
				currRoom = newEmployee.GetCurrentLocation()
			}
		// print success message only if there weren't any errors
			if ok {
				fmt.Println("\nSuccess! Currently at", office.GetLocationName(currRoom))
			}
		}
	}
	fmt.Println("****************Done for the day!******************************")
	displayEmplStatus(listOfEmployees)
}

// print all employees in the list with their location data
func displayEmplStatus(l []*office.Employee) {
	fmt.Println("Employee status:")
	for i, e := range l {
		// fmt.Printf("%#v\n", e) // useful code
		fmt.Printf("Employee %d (%s) is at the %s\n", i+1, office.GetEmployeeTitle(*e), office.GetLocationName((*e).GetCurrentLocation()))
	}
}

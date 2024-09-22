package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type model struct {
	name   string
	number int
}

func (m *model) setvalues(name string, number int) string {
	m.name = name
	m.number = number
	return fmt.Sprintf("Values added successfully")
}

func getdata(models []model) {
	for _, m := range models {
		fmt.Println(m.name, ": ", m.number)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	values := []model{}
	var choice string
	for {
		fmt.Println("Add more data: {yes/y/enter to exit}")
		fmt.Scanln(&choice)
		choice = strings.TrimSpace(choice)
		if choice == "yes" || choice == "y" {
			name, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			name = strings.TrimSpace(name)

			number, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
			number = strings.TrimSpace(number)
			num, err := strconv.Atoi(number)
			if err != nil {
				fmt.Printf("Error changing string to int: %v", err)
			}
			var m model
			m.setvalues(name, num)
			values = append(values, m)
		} else {
			break
		}
	}
	getdata(values)
}

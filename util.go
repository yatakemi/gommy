package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

// FailOnError function
func FailOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

// Exists function is checking a file exists
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// Question function
func Question(q string) bool {
	result := true
	fmt.Print(color.GreenString(q))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		i := scanner.Text()

		if i == "Y" || i == "y" || i == "" {
			break
		} else if i == "N" || i == "n" {
			result = false
			break
		} else {
			fmt.Println(color.RedString("Please answer Y or N"))
			fmt.Print(color.GreenString(q))
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fatih/color"
)

var (
	listSlice       []string
	packagesIDSlice []string
)

func main() {
	clearScreen()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Println(blue("Current time is:", time.Now()))
	fmt.Println()
	fmt.Println(yellow("Update list first."))
	fmt.Println()

	for {
		printPackageList()
		println()
		fmt.Println(yellow("1 Update List"))
		fmt.Println(yellow("2 Update Package"))
		fmt.Println()
		num, _ := readCommand(yellow("Enter cmd to continue, x to exit: "))

		// Switch statement to handle different cases
		switch num {
		case 1:
			fmt.Println(yellow("Updating list..."))
			fmt.Println()
			saveToFile()
			readFile()
			// println()
			// printPackageList()
			// println()
		case 2:
			// clearScreen()
			fmt.Println(yellow("Update Package..."))
			fmt.Println()
			printPackageList()
			fmt.Println()
			num, _ := readCommand(yellow("Enter package nr to update, x to exit, r to menu: "))
			fmt.Println()
			if num >= 0 {
				executeUpdateCommand(num)
				// listSlice = nil
				// packagesIDSlice = nil
				saveToFile()
				readFile()
			}
		default:
			fmt.Println(yellow("Unknown command."))
		}
	}
}

func clearScreen() {
	// For Windows
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		// For other systems (Unix-like)
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// func waitForString() {
// 	fmt.Print("Press Enter to continue...")
// 	var userInput string
// 	for {
// 		fmt.Scanln(&userInput)
// 		break
// 	}
// }

func startsWithLetter(s string) bool {
	if len(s) == 0 {
		return false
	}
	return unicode.IsLetter(rune(s[0]))
}

func saveToFile() {
	yellow := color.New(color.FgYellow).SprintFunc()
	file, err := os.Create("list.txt")
	if err != nil {
		fmt.Println(yellow("Error creating file:", err))
		return
	}
	defer file.Close()

	// Execute external command
	cmd := exec.Command("winget", "update")
	cmd.Stdout = file
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println(yellow("Error executing command:", err))
	}
}

func readFile() {
	yellow := color.New(color.FgYellow).SprintFunc()
	listSlice = nil
	packagesIDSlice = nil

	// Open the text file
	file, err := os.Open("list.txt")
	if err != nil {
		fmt.Println(yellow("Error opening file:", err))
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Read lines until the end of the file
	i := 0
	for scanner.Scan() {
		// Process each line here
		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		if !startsWithLetter(line) {
			continue
		}

		//hacks
		if strings.Contains(line, "Name    Id              Version  Available Source") ||
			strings.Contains(line, "The following packages have an upgrade available, but require explicit targeting for upgrade:") {
			listSlice = append(listSlice, line)
			packagesIDSlice = append(packagesIDSlice, "")
			i++
			continue
		}

		// Split the string into words based on whitespace characters
		words := strings.Fields(line)

		//hacks
		if len(words) > 4 {
			word4 := strings.Join(words[len(words)-4:len(words)-3], "")

			if len(word4) > 2 {
				listSlice = append(listSlice, line)
				packagesIDSlice = append(packagesIDSlice, word4)
			} else {
				word4 = strings.Join(words[len(words)-5:len(words)-4], "")
				listSlice = append(listSlice, line)
				packagesIDSlice = append(packagesIDSlice, word4)
			}
			i++
		}
	}

	// Check for any errors encountered during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println(yellow("Error reading file:", err))
	}
}

func printPackageList() {
	yellow := color.New(color.FgYellow).SprintFunc()

	//display list + packages id
	for i := range listSlice {
		//hacks
		idx := fmt.Sprintf("%*d", 3, i)
		if packagesIDSlice[i] == "" {
			fmt.Println("   ", listSlice[i], yellow(packagesIDSlice[i]))
		} else {
			fmt.Println(yellow(idx), listSlice[i], yellow(packagesIDSlice[i]))
		}
	}
}

func readCommand(str string) (int, string) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Print(str)
	var userInput string
	for {
		fmt.Scanln(&userInput)
		fmt.Println()
		if userInput == "x" {
			fmt.Println(yellow("Exiting program."))

			// Exit the program with status code 0
			os.Exit(0)
		}
		if userInput == "r" {
			fmt.Println(yellow("Return to menu."))
			return -1, ""
		}
		if startsWithLetter(userInput) {
			return -2, ""
		}

		// Convert string to integer
		num, err := strconv.Atoi(userInput)
		if err != nil {
			// Handle error if conversion fails
			fmt.Println(yellow("Error:", err))
		}
		return num, userInput
	}
}

func executeUpdateCommand(num int) {
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Println(yellow("Executing: winget update", packagesIDSlice[num]))

	//PROD
	cmd := exec.Command("winget", "update", packagesIDSlice[num])

	//TEST
	// fmt.Println(yellow("TEST CMD is \"winget update\" hence the following weird table!!!"))
	// cmd := exec.Command("winget", "update")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println(yellow("Error executing command:", err))
	}
}

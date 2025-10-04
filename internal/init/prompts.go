package initcmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// promptString asks the user for a string input with a default value
func promptString(question, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	if defaultValue != "" {
		fmt.Printf("? %s [%s]: ", question, defaultValue)
	} else {
		fmt.Printf("? %s: ", question)
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}

	return input, nil
}

// promptInt asks the user for an integer input with a default value
func promptInt(question string, defaultValue int) (int, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("? %s [%d]: ", question, defaultValue)

	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(input)
	if err != nil {
		return 0, fmt.Errorf("invalid number: %s", input)
	}

	return value, nil
}

// promptSelect asks the user to select from a list of options
func promptSelect(question string, options []string, defaultValue string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("? %s\n", question)
	for i, option := range options {
		prefix := " "
		if option == defaultValue {
			prefix = ">"
		}
		fmt.Printf("  %s %d) %s\n", prefix, i+1, option)
	}
	fmt.Printf("Select [1-%d] (default: %s): ", len(options), defaultValue)

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue, nil
	}

	// Try to parse as number
	index, err := strconv.Atoi(input)
	if err == nil {
		if index < 1 || index > len(options) {
			return "", fmt.Errorf("invalid selection: %d", index)
		}
		return options[index-1], nil
	}

	// Try to match as string
	for _, option := range options {
		if strings.EqualFold(input, option) {
			return option, nil
		}
	}

	return "", fmt.Errorf("invalid selection: %s", input)
}

// promptConfirm asks the user for a yes/no confirmation
func promptConfirm(question string, defaultValue bool) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	defaultStr := "y/N"
	if defaultValue {
		defaultStr = "Y/n"
	}

	fmt.Printf("? %s [%s]: ", question, defaultStr)

	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defaultValue, nil
	}

	switch input {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		return false, fmt.Errorf("invalid input: %s", input)
	}
}

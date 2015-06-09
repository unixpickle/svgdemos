package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type PathCommand struct {
	name      string
	arguments []float64
}

func commandFromStringArgs(name string, arguments []string) (*PathCommand, error) {
	res := PathCommand{name, []float64{}}
	for _, s := range arguments {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		res.arguments = append(res.arguments, f)
	}
	return &res, nil
}

func main() {
	fmt.Println("Enter path data, then deliver an EOF:")
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read data:", err)
		os.Exit(1)
	}
	parsed, err := parsePath(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse:", err)
		os.Exit(1)
	}
	a, b, c, d := pathBounds(parsed)
	fmt.Println("viewBox:", a, b, c, d)
}

func parsePath(data []byte) ([]PathCommand, error) {
	commands := []PathCommand{}
	args := []string{}
	var command string
	var currentArgument string
	for _, r := range []rune(string(data)) {
		if unicode.IsLetter(r) {
			if command != "" {
				if currentArgument != "" {
					args = append(args, currentArgument)
				}
				if c, err := commandFromStringArgs(command, args); err != nil {
					return nil, err
				} else {
					commands = append(commands, *c)
				}
			}
			command = string(r)
			currentArgument = ""
			args = []string{}
		} else if unicode.IsDigit(r) {
			currentArgument += string(r)
		} else if r == '.' {
			if strings.Contains(currentArgument, ".") {
				return nil, errors.New("extraneous decimal point")
			}
			currentArgument += "."
		} else if r == '-' {
			if currentArgument == "" {
				currentArgument = "-"
			} else {
				// The minus sign can be used as a separator.
				args = append(args, currentArgument)
				currentArgument = "-"
			}
		} else {
			if currentArgument != "" {
				args = append(args, currentArgument)
				currentArgument = ""
			}
		}
	}
	if command != "" {
		if currentArgument != "" {
			args = append(args, currentArgument)
		}
		if c, err := commandFromStringArgs(command, args); err != nil {
			return nil, err
		} else {
			commands = append(commands, *c)
		}
	}
	return commands, nil
}

func pathBounds(commands []PathCommand) (minX, minY, maxX, maxY float64) {
	minX = math.MaxFloat64
	minY = math.MaxFloat64
	maxX = -math.MaxFloat64
	maxY = -math.MaxFloat64
	
	// TODO: implement this.
	
	return
}

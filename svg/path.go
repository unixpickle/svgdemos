package svg

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type Bounder interface {
	Bounds() Rect
}

type PathCmd struct {
	Name string
	Args []float64
}

func (c PathCmd) Clone() PathCmd {
	res := PathCmd{c.Name, make([]float64, len(c.Args))}
	copy(res.Args, c.Args)
	return res
}

type Path []PathCmd

func ParsePath(s string) (Path, error) {
	path := Path{}

	args := []float64{}
	name := ""
	argStr := ""

	// NOTE: we add a "z" command to the path so that the command before it gets
	// pushed to the resulting path.
	for _, r := range s + "z" {
		isArg := unicode.IsDigit(r) || r == '.'
		if !isArg {
			if argStr != "" {
				if num, err := strconv.ParseFloat(argStr, 0); err != nil {
					return nil, err
				} else {
					args = append(args, num)
				}
				argStr = ""
			}
		}
		if unicode.IsLetter(r) {
			if name != "" {
				path = append(path, PathCmd{name, args})
				args = []float64{}
			}
			name = string(r)
		} else if isArg || r == '-' {
			argStr += string(r)
		}
	}
	if err := path.Validate(); err != nil {
		return nil, err
	} else {
		return path, nil
	}
}

// Absolute generates a path which only uses absolute commands.
func (p Path) Absolute() Path {
	currentPoint := Point{0, 0}
	subpathStart := Point{0, 0}
	res := make(Path, len(p))
	for i, cmd := range p {
		argCount := len(cmd.Args)
		if cmd.Name == "z" || cmd.Name == strings.ToUpper(cmd.Name) {
			res[i] = cmd.Clone()
		}
		switch cmd.Name {
		case "M":
			subpathStart := Point{cmd.Args[0], cmd.Args[1]}
			fallthrough
		case "L":
			currentPoint := Point{cmd.Args[argCount-2], cmd.Args[argCount-1]}
		case "m":
			subpathStart := Point{cmd.Args[0] + currentPoint.X, cmd.Args[1] +
				currentPoint.Y}
			fallthrough
		case "l":
			absCommand := PathCmd{strings.ToUpper(cmd.Name), []float64{}}
			for i := 0; i < argCount; i += 2 {
				currentPoint.X += cmd.Args[i]
				currentPoint.Y += cmd.Args[i+1]
				absCommand.Args = append(absCommand.Args, currentPoint.X,
					currentPoint.Y)
			}
			res[i] = absCommand
		case "Z", "z":
			currentPoint := subpathStart
		case "H":
			currentPoint.X = cmd.Args[argCount-1]
		case "h":
			absCommand := PathCmd{"H", []float64{}}
			for _, x := range cmd.Args {
				currentPoint.X += x
				absCommand.Args = append(absCommand.Args, currentPoint.X)
			}
			res[i] = absCommand
		case "V":
			currentPoint.Y = cmd.Args[argCount-1]
		case "v":
			absCommand := PathCmd{"V", []float64{}}
			for _, y := range cmd.Args {
				currentPoint.Y += y
				absCommand.Args = append(absCommand.Args, currentPoint.Y)
			}
			res[i] = absCommand
		case "C":
			currentPoint.X = cmd.Args[argCount-2]
			currentPoint.Y = cmd.Args[argCount-1]
		case "c":
			absCommand := PathCmd{"C", []float64{}}
			for i := 0; i < argCount; i += 6 {
				absCommand.Args = append(absCommand.Args,
					cmd.Args[i]+currentPoint.X, cmd.Args[i+1]+currentPoint.Y,
					cmd.Args[i+2]+currentPoint.X, cmd.Args[i+3]+currentPoint.Y,
					cmd.Args[i+4]+currentPoint.X, cmd.Args[i+5]+currentPoint.Y)
				currentPoint.X += cmd.Args[i+4]
				currentPoint.Y += cmd.Args[i+5]
			}
			res[i] = absCommand
		case "S":
			currentPoint.X = cmd.Args[argCount-2]
			currentPoint.Y = cmd.Args[argCount-1]
		case "s":
			absCommand := PathCmd{"S", []float64{}}
			for i := 0; i < argCount; i += 4 {
				absCommand.Args = append(absCommand.Args,
					cmd.Args[i]+currentPoint.X, cmd.Args[i+1]+currentPoint.Y,
					cmd.Args[i+2]+currentPoint.X, cmd.Args[i+3]+currentPoint.Y)
				currentPoint.X += cmd.Args[i+2]
				currentPoint.Y += cmd.Args[i+3]
			}
			res[i] = absCommand

		// TODO: do these commands
		case "Q":
		case "q":
		case "T":
		case "t":
		case "A":
		case "a":
		}
	}
	return res
}

// Validate makes sure that the path has valid commands and arguments. If not,
// it returns an error describing the problem.
func (p Path) Validate() error {
	argCounts := map[string]int{"m": 2, "z": 0, "l": 2, "h": 1, "v": 1, "c": 6,
		"s": 4, "q": 4, "t": 2, "a": 7}
	for _, cmd := range p {
		lowerName := strings.ToLower(cmd.Name)
		count, ok := argCounts[lowerName]
		if !ok {
			return errors.New("unknown command: " + cmd.Name)
		} else if lowerName == "z" && len(cmd.Args) != 0 {
			return errors.New(cmd.Name + " command takes no arguments")
		} else if len(cmd.Args) < count {
			return errors.New("not enough arguments to " + cmd.Name)
		} else if len(cmd.Args)%count != 0 {
			return errors.New("invalid number of arguments to " + cmd.Name)
		}
	}
	return nil
}

package svg

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type PathShape interface {
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

func (c PathCmd) Equals(x PathCmd) bool {
	if c.Name != x.Name || len(c.Args) != len(x.Args) {
		return false
	}
	for i, arg := range c.Args {
		if arg != x.Args[i] {
			return false
		}
	}
	return true
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
			if name == "" {
				return nil, errors.New("argument before first command name")
			}
			argStr += string(r)
		}
	}
	if err := path.Validate(); err != nil {
		return nil, err
	} else {
		return path, nil
	}
}

// Absolute generates a path which only uses absolute commands. A path must be
// validated before it can be made absolute.
func (p Path) Absolute() Path {
	if err := p.Validate(); err != nil {
		panic("path is invalid: " + err.Error())
	}

	currentPoint := Point{0, 0}
	subpathStart := Point{0, 0}
	res := make(Path, len(p))
	for i, cmd := range p {
		argCount := len(cmd.Args)
		if cmd.Name == "z" || cmd.Name == strings.ToUpper(cmd.Name) {
			res[i] = cmd.Clone()
			res[i].Name = strings.ToUpper(cmd.Name)
		}
		switch cmd.Name {
		case "M":
			subpathStart.X = cmd.Args[0]
			subpathStart.Y = cmd.Args[1]
			fallthrough
		case "L", "C", "S", "Q", "T", "A":
			currentPoint.X = cmd.Args[argCount-2]
			currentPoint.Y = cmd.Args[argCount-1]
		case "m":
			subpathStart.X = cmd.Args[0] + currentPoint.X
			subpathStart.Y = cmd.Args[1] + currentPoint.Y
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
			currentPoint = subpathStart
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
		case "q":
			absCommand := PathCmd{"Q", []float64{}}
			for i := 0; i < argCount; i += 4 {
				absCommand.Args = append(absCommand.Args,
					cmd.Args[i]+currentPoint.X, cmd.Args[i+1]+currentPoint.Y,
					cmd.Args[i+2]+currentPoint.X, cmd.Args[i+3]+currentPoint.Y)
				currentPoint.X += cmd.Args[i+2]
				currentPoint.Y += cmd.Args[i+3]
			}
			res[i] = absCommand
		case "t":
			absCommand := PathCmd{"T", []float64{}}
			for i := 0; i < argCount; i += 2 {
				absCommand.Args = append(absCommand.Args,
					cmd.Args[i]+currentPoint.X, cmd.Args[i+1]+currentPoint.Y)
				currentPoint.X += cmd.Args[i]
				currentPoint.Y += cmd.Args[i+1]
			}
			res[i] = absCommand
		case "a":
			absCommand := PathCmd{"A", []float64{}}
			for i := 0; i < argCount; i += 7 {
				absCommand.Args = append(absCommand.Args, cmd.Args[i:i+5]...)
				currentPoint.X += cmd.Args[i+5]
				currentPoint.Y += cmd.Args[i+6]
				absCommand.Args = append(absCommand.Args, currentPoint.X,
					currentPoint.Y)
			}
			res[i] = absCommand
		}
	}
	return res
}

// Normalize performs a number of transformations to a path to make it easier to
// process and read. A path must be validated before it can be normalized.
//
// A normalized path has no relative commands; all commands are absolute. When a
// command is called multiple times in a row, a normalized path has a separate
// PathCmd element for each of these calls. For example, a normalized path would
// represent "L 10,10 20,0" as {PathCmd{"L", {10, 10}}, PathCmd{"L", {20, 0}}}.
// Shorthand commands like "H", "V", "S" and "T" are converted into their longer
// equivalents "L", "C" and "Q".
func (p Path) Normalize() Path {
	if err := p.Validate(); err != nil {
		panic("path is invalid: " + err.Error())
	}

	multicalls := p.Absolute().SplitMulticalls()
	res := make(Path, 0, len(multicalls))

	// Convert shorthand commands into longhand.
	currentPoint := Point{0, 0}
	subpathStart := Point{0, 0}
	for _, cmd := range multicalls {
		argCount := len(cmd.Args)
		switch cmd.Name {
		case "H":
			currentPoint.X = cmd.Args[0]
			line := PathCmd{"L", []float64{currentPoint.X, currentPoint.Y}}
			res = append(res, line)
		case "V":
			currentPoint.Y = cmd.Args[0]
			line := PathCmd{"L", []float64{currentPoint.X, currentPoint.Y}}
			res = append(res, line)
		case "S":
			lastControlPoint := currentPoint
			if len(res) > 0 && res[len(res)-1].Name == "C" {
				lastCmd := res[len(res)-1]
				lastControlPoint = Point{lastCmd.Args[2], lastCmd.Args[3]}
			}
			reflected := []float64{currentPoint.X*2 - lastControlPoint.X,
				currentPoint.Y*2 - lastControlPoint.Y}
			res = append(res, PathCmd{"C", append(reflected, cmd.Args...)})
		case "T":
			lastControlPoint := currentPoint
			if len(res) > 0 && res[len(res)-1].Name == "Q" {
				lastCmd := res[len(res)-1]
				lastControlPoint = Point{lastCmd.Args[0], lastCmd.Args[1]}
			}
			reflected := []float64{currentPoint.X*2 - lastControlPoint.X,
				currentPoint.Y*2 - lastControlPoint.Y}
			res = append(res, PathCmd{"Q", append(reflected, cmd.Args...)})
		default:
			res = append(res, cmd)
			if cmd.Name == "Z" {
				currentPoint = subpathStart
			} else if cmd.Name == "M" {
				subpathStart = Point{cmd.Args[0], cmd.Args[1]}
			}
		}
		if len(cmd.Args) >= 2 {
			currentPoint = Point{cmd.Args[argCount-2], cmd.Args[argCount-1]}
		}
	}

	return res
}

// SplitMulticalls separates consecutive calls to the same command into separate
// PathCmd objects in a path. It will also split up multicalls to the moveto
// command, so thinks like "M a,b,c,d" are turned into "M a,b L c,d".
//
// A path must be validated before SplitMulticalls can work.
func (p Path) SplitMulticalls() Path {
	if err := p.Validate(); err != nil {
		panic("path is invalid: " + err.Error())
	}

	res := make(Path, 0, len(p))
	argCounts := map[string]int{"m": 2, "z": 0, "l": 2, "h": 1, "v": 1, "c": 6,
		"s": 4, "q": 4, "t": 2, "a": 7}
	for _, cmd := range p {
		if cmd.Name == "z" || cmd.Name == "Z" {
			res = append(res, cmd.Clone())
			continue
		}
		argCount := argCounts[strings.ToLower(cmd.Name)]
		for i := 0; i < len(cmd.Args); i += argCount {
			subArgs := cmd.Args[i : i+argCount]
			argCopy := make([]float64, argCount)
			copy(argCopy, subArgs)
			name := cmd.Name
			if name == "M" && i > 0 {
				name = "L"
			} else if name == "m" && i > 0 {
				name = "l"
			}
			res = append(res, PathCmd{name, argCopy})
		}
	}

	return res
}

// String exports the path as a string which can be used in an SVG file.
func (p Path) String() string {
	var buffer bytes.Buffer
	for _, cmd := range p {
		buffer.WriteString(cmd.Name)
		for i, arg := range cmd.Args {
			argStr := strconv.FormatFloat(arg, 'f', -1, 64)
			if i > 0 && argStr[0] != '-' {
				buffer.WriteRune(' ')
			}
			buffer.WriteString(argStr)
		}
	}
	return buffer.String()
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
		} else if lowerName == "z" {
			if len(cmd.Args) != 0 {
				return errors.New(cmd.Name + " command takes no arguments")
			}
		} else if len(cmd.Args) < count {
			return errors.New("not enough arguments to " + cmd.Name)
		} else if len(cmd.Args)%count != 0 {
			return errors.New("invalid number of arguments to " + cmd.Name)
		}
	}
	return nil
}

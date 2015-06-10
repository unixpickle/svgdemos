package svg

import (
	"strconv"
	"unicode"
)

type Bounder interface {
	Bounds() Rect
}

type PathCmd struct {
	Name string
	Args []float64
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
	return path, nil
}

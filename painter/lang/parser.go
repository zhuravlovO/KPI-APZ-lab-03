package lang

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"image"

	"github.com/zhuravlovO/KPI-APZ-lab-03/painter"
)

// Parser вміє прочитати дані з вхідного io.Reader та повернути список операцій.
type Parser struct{}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)
	const (
		screenWidth  = 800
		screenHeight = 800
	)

	for scanner.Scan() {
		commandLine := scanner.Text()
		fmt.Println("Parser got line:", commandLine)

		fields := strings.Fields(commandLine)
		if len(fields) == 0 {
			continue
		}

		cmd := fields[0]
		args := fields[1:]

		switch cmd {
		case "white":
			res = append(res, painter.WhiteFill)
		case "green":
			res = append(res, painter.GreenFill)
		case "update":
			res = append(res, painter.UpdateOp)
		case "reset":
			res = append(res, painter.ResetOp)
		case "bgrect":
			if len(args) != 4 {
				continue
			}
			coords := make([]float64, 4)
			for i, arg := range args {
				val, err := strconv.ParseFloat(arg, 64)
				if err != nil {
					continue
				}
				coords[i] = val
			}
			rect := image.Rect(
				int(coords[0]*screenWidth), int(coords[1]*screenHeight),
				int(coords[2]*screenWidth), int(coords[3]*screenHeight),
			)
			res = append(res, painter.BgRectOp{Rect: rect})
		case "figure":
			if len(args) != 2 {
				continue
			}
			x, _ := strconv.ParseFloat(args[0], 64)
			y, _ := strconv.ParseFloat(args[1], 64)
			res = append(res, painter.FigureOp{X: int(x * screenWidth), Y: int(y * screenHeight)})
		case "move":
			if len(args) != 2 {
				continue
			}
			x, _ := strconv.ParseFloat(args[0], 64)
			y, _ := strconv.ParseFloat(args[1], 64)
			res = append(res, painter.MoveOp{X: int(x * screenWidth), Y: int(y * screenHeight)})
		}
	}

	return res, nil
}

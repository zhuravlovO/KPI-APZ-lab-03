package painter

import (
	"image"
	"image/color"
	"image/draw"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	Do(l *Loop) bool
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(l *Loop) (ready bool) {
	for _, o := range ol {
		ready = o.Do(l) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(l *Loop) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(l *Loop) bool

func (f OperationFunc) Do(l *Loop) bool {
	return f(l)
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(l *Loop) bool {
	l.next.Fill(l.next.Bounds(), color.White, draw.Src)
	return true
}

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(l *Loop) bool {
	green := color.RGBA{G: 0xff, A: 0xff}
	l.next.Fill(l.next.Bounds(), green, draw.Src)
	return true
}

type MoveFigureOp struct {
	X, Y int
}

func (op MoveFigureOp) Do(l *Loop) bool {
	l.figureX = op.X
	l.figureY = op.Y
	l.next.Fill(l.next.Bounds(), color.White, draw.Src)

	blue := color.RGBA{B: 0xff, A: 0xff}
	verticalRect := image.Rect(l.figureX-50, l.figureY-150, l.figureX+50, l.figureY+150)
	horizontalRect := image.Rect(l.figureX-150, l.figureY-50, l.figureX+150, l.figureY+50)

	l.next.Fill(verticalRect, blue, draw.Src)
	l.next.Fill(horizontalRect, blue, draw.Src)
	return true
}

package painter

import (
	"image"
	"image/color"
	"image/draw"
)

// Operation — це інтерфейс для будь-якої дії, що змінює стан сцени.
type Operation interface {
	Do(l *Loop) bool
}

type OperationList []Operation

func (ol OperationList) Do(l *Loop) (ready bool) {
	for _, o := range ol {
		ready = o.Do(l) || ready
	}
	return
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(l *Loop) bool {
	l.next.Fill(l.next.Bounds(), l.BgColor, draw.Src)
	if l.BgRect != nil {
		l.next.Fill(*l.BgRect, color.Black, draw.Src)
	}
	blue := color.RGBA{B: 0xff, A: 0xff}
	for _, fig := range l.Figures {
		verticalRect := image.Rect(fig.X+50, fig.Y-200, fig.X+150, fig.Y+200)
		horizontalRect := image.Rect(fig.X-200, fig.Y-50, fig.X+50, fig.Y+50)
		l.next.Fill(verticalRect, blue, draw.Src)
		l.next.Fill(horizontalRect, blue, draw.Src)
	}
	return true
}

type OperationFunc func(l *Loop) bool

func (f OperationFunc) Do(l *Loop) bool {
	return f(l)
}

var (
	WhiteFill = OperationFunc(func(l *Loop) bool {
		l.BgColor = color.White
		return false
	})
	GreenFill = OperationFunc(func(l *Loop) bool {
		l.BgColor = color.RGBA{G: 0xff, A: 0xff}
		return false
	})
)
var ResetOp = OperationFunc(func(l *Loop) bool {
	l.BgColor = color.Black
	l.BgRect = nil
	l.Figures = make([]*Figure, 0)
	return false
})

type BgRectOp struct {
	Rect image.Rectangle
}

func (op BgRectOp) Do(l *Loop) bool {
	l.BgRect = &op.Rect
	return false
}

type FigureOp struct {
	X, Y int
}

func (op FigureOp) Do(l *Loop) bool {
	l.Figures = append(l.Figures, &Figure{X: op.X, Y: op.Y})
	return false
}

type MoveOp struct {
	X, Y int
}

func (op MoveOp) Do(l *Loop) bool {
	for _, fig := range l.Figures {
		fig.X += op.X
		fig.Y += op.Y
	}
	return false
}

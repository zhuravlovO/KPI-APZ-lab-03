package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/mouse"
)

type Figure struct {
	X, Y int
}
type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver
	next     screen.Texture
	prev     screen.Texture

	Mq chan Operation

	stop    chan struct{}
	stopReq bool

	BgColor color.Color
	BgRect  *image.Rectangle
	Figures []*Figure
}

// Start запускає цикл подій.
func (l *Loop) Start(s screen.Screen) {
	size := image.Point{X: 800, Y: 800}
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)
	l.stop = make(chan struct{})

	l.BgColor = color.Black
	l.Figures = make([]*Figure, 0)

	go func() {
		for op := range l.Mq {
			if op.Do(l) {
				l.Receiver.Update(l.next)
			}
		}
		close(l.stop)
	}()
	l.Post(FigureOp{X: 400, Y: 400})
	l.Post(UpdateOp)
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	l.Mq <- op
}
func (l *Loop) HandleMouse(e mouse.Event) {
	if e.Button == mouse.ButtonRight {
		ops := OperationList{
			OperationFunc(func(l *Loop) bool {
				l.Figures = make([]*Figure, 0)
				return false
			}),
			FigureOp{X: int(e.X), Y: int(e.Y)},
			UpdateOp,
		}
		l.Post(ops)
	}
}
func (l *Loop) StopAndWait() {
	if !l.stopReq {
		l.stopReq = true
		close(l.Mq)
	}
	<-l.stop
}

package painter

import (
	"image"

	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/mouse"
)

type Receiver interface {
	Update(t screen.Texture)
}

type Loop struct {
	Receiver Receiver
	next     screen.Texture
	prev     screen.Texture
	Mq       chan Operation

	stop    chan struct{}
	stopReq bool

	figureX, figureY int
}

func (l *Loop) Start(s screen.Screen) {
	size := image.Point{X: 800, Y: 800}
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.stop = make(chan struct{})

	l.figureX = 400
	l.figureY = 400

	go func() {
		for op := range l.Mq {
			if op.Do(l) {
				l.Receiver.Update(l.next)
			}
		}
		close(l.stop)
	}()
}

func (l *Loop) Post(op Operation) {
	l.Mq <- op
}

func (l *Loop) HandleMouse(e mouse.Event) {
	if e.Button == mouse.ButtonRight {
		op := MoveFigureOp{
			X: int(e.X),
			Y: int(e.Y),
		}
		l.Post(op)
	}
}

func (l *Loop) StopAndWait() {
	if !l.stopReq {
		l.stopReq = true
		close(l.Mq)
	}
	<-l.stop
}

package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
	"time"

	"golang.org/x/exp/shiny/screen"
)

func TestLoopPost(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)

	l.Mq = make(chan Operation, 100)
	l.Receiver = &tr

	l.Start(mockScreen{})

	l.Post(OperationFunc(WhiteFill))
	l.Post(OperationFunc(GreenFill))

	l.Post(UpdateOp)

	time.Sleep(100 * time.Millisecond)

	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", tr.lastTexture)
	}
	finalBgColor := color.RGBA{G: 0xff, A: 0xff}
	if mt.bgColor != finalBgColor {
		t.Errorf("Unexpected background color: got %v, want %v", mt.bgColor, finalBgColor)
	}
}

type mockTexture struct {
	bgColor color.Color
}

func (m *mockTexture) Release()                                                     {}
func (m *mockTexture) Size() image.Point                                            { return image.Point{X: 800, Y: 800} }
func (m *mockTexture) Bounds() image.Rectangle                                      { return image.Rectangle{Max: m.Size()} }
func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	if dr == m.Bounds() {
		m.bgColor = src
	}
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

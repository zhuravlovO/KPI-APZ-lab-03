package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

func TestLoopPost(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr
	l.Mq = make(chan Operation)

	var testOps []string
	l.Start(mockScreen{})
	l.Post(logOp(t, "do white fill", OperationFunc(WhiteFill)))
	l.Post(logOp(t, "do green fill", OperationFunc(GreenFill)))
	l.Post(UpdateOp)

	for i := 0; i < 3; i++ {
		go l.Post(logOp(t, "do green fill", OperationFunc(GreenFill)))
	}

	l.Post(OperationFunc(func(l *Loop) bool {
		testOps = append(testOps, "op 1")
		return false
	}))
	l.Post(OperationFunc(func(l *Loop) bool {
		testOps = append(testOps, "op 2")
		return false
	}))

	l.Post(OperationFunc(func(l *Loop) bool {
		testOps = append(testOps, "op 3")
		return true // Ця операція викличе оновлення
	}))

	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", tr.lastTexture)
	}
	if mt.Colors[0] != color.White {
		t.Error("First color is not white:", mt.Colors)
	}
	if len(mt.Colors) != 2 { // WhiteFill + GreenFill
		t.Error("Unexpected size of colors:", mt.Colors)
	}

	if !reflect.DeepEqual([]string{"op 1", "op 2", "op 3"}, testOps) {
		t.Error("Bad order:", testOps)
	}
}

func logOp(t *testing.T, msg string, op OperationFunc) OperationFunc {
	return func(l *Loop) bool {
		t.Log(msg)
		return op(l)
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

type mockTexture struct {
	Colors []color.Color
	size   image.Point
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point {
	return m.size
}
func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.size}
}
func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {
	panic("implement me")
}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}

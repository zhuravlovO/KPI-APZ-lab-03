package main

import (
	"net/http"

	"github.com/zhuravlovO/KPI-APZ-lab-03/painter"
	"github.com/zhuravlovO/KPI-APZ-lab-03/painter/lang"
	"github.com/zhuravlovO/KPI-APZ-lab-03/ui"
)

func main() {
	var (
		pv     ui.Visualizer
		opLoop painter.Loop
		parser lang.Parser
	)

	opLoop.Mq = make(chan painter.Operation)

	//pv.Debug = true
	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	pv.OnMouseEvent = opLoop.HandleMouse
	opLoop.Receiver = &pv

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, &parser))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	pv.Main()
	opLoop.StopAndWait()
}

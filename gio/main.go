package main

import (
	"gioui.org/ui"
	"gioui.org/ui/app"
	"gioui.org/ui/layout"
)

func main() {
	go func() {
		w := app.NewWindow()
		ops := new(ui.Ops)
		for e := range w.Events() {
			switch e := e.(type) {
			case app.UpdateEvent:
				ops.Reset()
				layout.RigidConstraints(e.Size)
				w.Update(ops)
			}
		}
	}()
	app.Main()
}

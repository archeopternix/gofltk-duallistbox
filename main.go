package main

import (
	"github.com/pwiecz/go-fltk"
)

func main() {

	window := fltk.NewWindow(400, 300)
	window.SetLabel("DualListBox Example")

	// Create a DualListBox at (20, 20), width 360, height 200
	dlb := NewDualListBox(window, 20, 20, 360, 200)
	dlb.SetLeftItems([]string{"Apple", "Banana", "Cherry", "Date"})
	dlb.SetRightItems([]string{})

	// Add an Exit button at the bottom
	exitBtn := fltk.NewButton(150, 240, 100, 40, "Exit")
	exitBtn.SetCallback(func() {
		fltk.Quit()
	})

	window.End()
	window.Show()
	fltk.Run()
}

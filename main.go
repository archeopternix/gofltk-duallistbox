package main

import (
	"os"

	"github.com/pwiecz/go-fltk"
)

func main() {

	window := fltk.NewWindow(400, 400)
	window.SetLabel("DualListBox Example")

	// Create a DualListBox at (20, 20), width 360, height 200
	dlb := NewDualListBox(window, 5, 5, window.W(), 200)
	dlb.SetLeftItems([]string{"Apple", "Banana", "Cherry", "Date"})
	dlb.SetRightItems([]string{})
	dlb.SetLeftTitle("Target")
	dlb.SetRightTitle("Source")
	window.Add(dlb)

	// Add an Exit button at the bottom
	exitBtn := fltk.NewButton(window.W()/2-50, window.H()-50, 100, 40, "Exit")
	exitBtn.SetCallback(func() {
		os.Exit(0)
	})

	window.End()
	window.Show()
	fltk.Run()
}

package main

import (
	"slices"
	"sort"

	"github.com/pwiecz/go-fltk"
)

// DualListBox is a widget for moving items between two lists.
type DualListBox struct {
	group                *fltk.Group
	leftBox, rightBox    *fltk.Box
	leftList, rightList  *fltk.MultiBrowser
	addButton, delButton *fltk.Button
	onClick              func()
	rightItems           []string // available filters (right side)
	leftItems            []string // used filters (left side)
	x, y, w, h           int
}

// NewDualListBox creates a dual list box UI element.
func NewDualListBox(parent *fltk.Window, x, y, w, h int) *DualListBox {
	group := fltk.NewGroup(x, y, w, h)
	group.Begin()

	boxH := 30
	btnW := 50
	btnH := 30
	btnX := w/2 - btnW/2
	btnY := y + boxH + (h-30)/2 // middle of the button area
	leftW := (w-btnW)/2 - 20
	rightW := (w-btnW)/2 - 20
	listH := h - 30
	// Left List
	leftBox := fltk.NewBox(fltk.NO_BOX, x, y, leftW, boxH, "Used Filters")
	leftBox.SetAlign(fltk.ALIGN_CENTER | fltk.ALIGN_INSIDE)
	leftList := fltk.NewMultiBrowser(x, y+boxH, leftW, listH, "")

	// Right List
	rightBox := fltk.NewBox(fltk.NO_BOX, x+w-rightW-10, y, rightW, boxH, "Available Filters")
	rightBox.SetAlign(fltk.ALIGN_CENTER | fltk.ALIGN_INSIDE)
	rightList := fltk.NewMultiBrowser(x+w-rightW-10, y+boxH, rightW, listH, "")

	addButton := fltk.NewButton(btnX, btnY-40, btnW, btnH, "<--")
	delButton := fltk.NewButton(btnX, btnY+10, btnW, btnH, "-->")

	dlb := &DualListBox{
		group:      group,
		leftBox:    leftBox,
		rightBox:   rightBox,
		leftList:   leftList,
		rightList:  rightList,
		addButton:  addButton,
		delButton:  delButton,
		rightItems: nil,
		leftItems:  nil,
		x:          x,
		y:          y,
		w:          w,
		h:          h,
	}

	addButton.SetCallback(func() {
		selected := rightList.Value()
		if selected > 0 {
			item := rightList.Text(selected)
			dlb.leftItems = append(dlb.leftItems, item)
			// Remove from rightItems
			index := slices.Index(dlb.rightItems, item)
			if index >= 0 {
				dlb.rightItems = slices.Delete(dlb.rightItems, index, index+1)
			}
			// Sort both lists and refresh
			sort.Strings(dlb.leftItems)
			sort.Strings(dlb.rightItems)
			dlb.Refresh()
			if dlb.onClick != nil {
				dlb.onClick()
			}
		}
	})
	delButton.SetCallback(func() {
		selected := leftList.Value()
		if selected > 0 {
			item := leftList.Text(selected)
			dlb.rightItems = append(dlb.rightItems, item)
			index := slices.Index(dlb.leftItems, item)
			if index >= 0 {
				dlb.leftItems = slices.Delete(dlb.leftItems, index, index+1)
			}
			// Sort both lists and refresh
			sort.Strings(dlb.leftItems)
			sort.Strings(dlb.rightItems)
			dlb.Refresh()
			if dlb.onClick != nil {
				dlb.onClick()
			}
		}
	})

	group.End()
	dlb.Refresh()
	return dlb
}

// SetLeftTitle sets the title of the left list box.
func (d *DualListBox) SetLeftTitle(title string) {
	if d.leftBox != nil {
		d.leftBox.SetLabel(title)
	}
}

// SetRightTitle sets the title of the right list box.
func (d *DualListBox) SetRightTitle(title string) {
	if d.rightBox != nil {
		d.rightBox.SetLabel(title)
	}
}

// RegisterClickHandler sets a callback called after add or delete button is pressed.
func (d *DualListBox) RegisterClickHandler(cb func()) {
	d.onClick = cb
}

// Refresh clears and repopulates both lists from the current rightItems and leftItems.
func (d *DualListBox) Refresh() {
	d.leftList.Clear()
	d.rightList.Clear()
	for _, f := range d.leftItems {
		d.leftList.Add(f)
	}
	for _, f := range d.rightItems {
		d.rightList.Add(f)
	}
}

// Resize moves and resizes the DualListBox and all its children.
func (d *DualListBox) Resize(x, y, w, h int) {
	d.x, d.y, d.w, d.h = x, y, w, h
	d.group.Resize(x, y, w, h)

	// Re-layout children (fixed widths for lists/buttons, adjust positions)
	btnW := 50
	btnH := 30
	btnX := x + w/2 - btnW/2
	btnY := h / 2 // middle of the button area
	leftW := (w+btnW)/2 - 20
	rightW := (w+btnW)/2 - 20
	listH := h - 30
	boxH := 30

	// Left
	d.leftList.Resize(x, y+boxH, leftW, listH)
	if d.leftBox != nil {
		d.leftBox.Resize(x, y, leftW, boxH)
	}

	// Right
	rightX := x + w - rightW
	d.rightList.Resize(rightX, y+boxH, rightW, listH)
	if d.rightBox != nil {
		d.rightBox.Resize(rightX, y, rightW, boxH)
	}

	// Buttons
	d.addButton.Resize(btnX, btnY-40, btnW, btnH)
	d.delButton.Resize(btnX, btnY+10, btnW, btnH)
}

// SetLeftItems replaces the items in the left list (used filters).
func (d *DualListBox) SetLeftItems(items []string) {
	d.leftList.Clear()
	sort.Strings(items)
	for _, item := range items {
		d.leftList.Add(item)
	}
	d.leftItems = slices.Clone(items)
}

// GetLeftItems returns all items in the left list (used filters).
func (d *DualListBox) GetLeftItems() []string {
	var items []string
	for i := 1; i <= d.leftList.Size(); i++ {
		items = append(items, d.leftList.Text(i))
	}
	return items
}

// SetRightItems replaces the items in the right list (available filters).
func (d *DualListBox) SetRightItems(items []string) {
	d.rightList.Clear()
	sort.Strings(items)
	for _, item := range items {
		d.rightList.Add(item)
	}
	d.rightItems = slices.Clone(items)
}

// GetRightItems returns all items in the right list (available filters).
func (d *DualListBox) GetRightItems() []string {
	var items []string
	for i := 1; i <= d.rightList.Size(); i++ {
		items = append(items, d.rightList.Text(i))
	}
	return items
}

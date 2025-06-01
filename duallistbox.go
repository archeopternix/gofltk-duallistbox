package main

import (
	"slices"

	"github.com/pwiecz/go-fltk"
)

// FilterList is a list of filter names.
type FilterList []string

func NewFilterList(filters ...string) FilterList {
	return FilterList(filters)
}

// DualListBox is a widget for moving items between two lists.
type DualListBox struct {
	group                *fltk.Group
	leftList, rightList  *fltk.MultiBrowser
	addButton, delButton *fltk.Button
	onClick              func()
	srcList              FilterList  // pointer to source (available) filters
	dstList              *FilterList // pointer to destination (used) filters
	x, y, w, h           int
}

// NewDualListBox creates a dual list box UI element.
// srcList: available filters (right side, pointer!)
// dstList: currently selected filters (left side, pointer!)
func NewDualListBox(parent *fltk.Window, x, y, w, h int, srcList FilterList, dstList *FilterList) *DualListBox {
	group := fltk.NewGroup(x, y, w, h)
	group.Begin()

	// Left List
	leftBox := fltk.NewBox(fltk.NO_BOX, x, y, 270, 30, "Used Filters")
	leftBox.SetAlign(fltk.ALIGN_CENTER | fltk.ALIGN_INSIDE)
	leftList := fltk.NewMultiBrowser(x, y+30, 270, h-30, "")
	for _, f := range *dstList {
		leftList.Add(f)
	}

	// Right List
	rightBox := fltk.NewBox(fltk.NO_BOX, x+w-270, y, 270, 30, "Available Filters")
	rightBox.SetAlign(fltk.ALIGN_CENTER | fltk.ALIGN_INSIDE)
	rightList := fltk.NewMultiBrowser(x+w-270, y+30, 270, h-30, "")
	for _, f := range srcList {
		if !slices.Contains(*dstList, f) {
			rightList.Add(f)
		}
	}

	btnX := x + w/2 - 50
	btnY := h / 2 // middle of the button area
	addButton := fltk.NewButton(btnX, btnY-40, 100, 30, "Add")
	delButton := fltk.NewButton(btnX, btnY+10, 100, 30, "Delete")

	dlb := &DualListBox{
		group:     group,
		leftList:  leftList,
		rightList: rightList,
		addButton: addButton,
		delButton: delButton,
		srcList:   srcList,
		dstList:   dstList,
		x:         x,
		y:         y,
		w:         w,
		h:         h,
	}

	addButton.SetCallback(func() {
		selected := rightList.Value()
		if selected > 0 {
			item := rightList.Text(selected)
			leftList.Add(item)
			*dlb.dstList = append(*dlb.dstList, item)
			rightList.Remove(selected)
			if dlb.onClick != nil {
				dlb.onClick()
			}
		}
	})
	delButton.SetCallback(func() {
		selected := leftList.Value()
		if selected > 0 {
			item := leftList.Text(selected)
			rightList.Add(item)
			index := slices.Index(*dlb.dstList, item)
			if index >= 0 {
				*dlb.dstList = slices.Delete(*dlb.dstList, index, index+1)
			}
			leftList.Remove(selected)
			if dlb.onClick != nil {
				dlb.onClick()
			}
		}
	})

	group.End()
	return dlb
}

// RegisterClickHandler sets a callback called after add or delete button is pressed.
func (d *DualListBox) RegisterClickHandler(cb func()) {
	d.onClick = cb
}

// Refresh clears and repopulates both lists from the current srcList and dstList.
func (d *DualListBox) Refresh() {
	d.leftList.Clear()
	d.rightList.Clear()
	for _, f := range *d.dstList {
		d.leftList.Add(f)
	}
	for _, f := range d.srcList {
		if !slices.Contains(*d.dstList, f) {
			d.rightList.Add(f)
		}
	}
}

// Resize moves and resizes the DualListBox and all its children.
func (d *DualListBox) Resize(x, y, w, h int) {
	d.x, d.y, d.w, d.h = x, y, w, h
	d.group.Resize(x, y, w, h)

	// Re-layout children (fixed widths for lists/buttons, adjust positions)
	leftW := 270
	rightW := 270
	listH := h - 30
	boxH := 30
	btnW := 100
	btnH := 30
	btnX := x + w/2 - btnW/2
	btnY := h / 2 // middle of the button area

	// Left
	d.leftList.Resize(x, y+boxH, leftW, listH)
	d.leftList.Parent().Resize(x, y, leftW, boxH) // leftBox (parent)

	// Right
	rightX := x + w - rightW
	d.rightList.Resize(rightX, y+boxH, rightW, listH)
	d.rightList.Parent().Resize(rightX, y, rightW, boxH) // rightBox (parent)

	// Buttons
	d.addButton.Resize(btnX, btnY-40, btnW, btnH)
	d.delButton.Resize(btnX, btnY+10, btnW, btnH)
}

// GetFilterList returns a pointer to the current used filters (left box).
func (d *DualListBox) GetFilterList() *FilterList {
	return d.dstList
}

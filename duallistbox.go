package duallistbox

import (
	"slices"
	"sort"

	"github.com/archeopternix/go-fltk"
)

// DualListBox is a composite widget for moving items between two lists.
// The left list represents "used" items and the right list "available" items.
// Users can move items between the lists using the arrow buttons.
// The widget is implemented by embedding a fltk.Group and managing its children.
type DualListBox struct {
	*fltk.Group
	leftBox, rightBox   *fltk.Box          // Label boxes above each list
	leftList, rightList *fltk.MultiBrowser // Multi-select lists for used/available items
	moveLeftButton      *fltk.Button       // Moves selected item from right to left
	moveRightButton     *fltk.Button       // Moves selected item from left to right
	onMoveRight         func()             // Handler called after a move to the right
	onMoveLeft          func()             // Handler called after a move to the left
	rightItems          []string           // Items shown on the right (available)
	leftItems           []string           // Items shown on the left (used)
	x, y, w, h          int                // Widget geometry
}

// NewDualListBox creates and returns a new DualListBox widget.
//
// x, y:   The top-left coordinate of the widget
// w, h:   The size of the widget
func NewDualListBox(x, y, w, h int) *DualListBox {
	group := fltk.NewGroup(x, y, w, h)
	group.Begin()

	boxH := 30
	btnW := 50
	btnH := 30
	btnX := w/2 - btnW/2
	btnY := y + boxH + (h-30)/2 // Vertically center buttons between the lists
	leftW := (w-btnW)/2 - 20
	rightW := (w-btnW)/2 - 20
	listH := h - 30

	// Left List: Used Items
	leftBox := fltk.NewBox(fltk.NO_BOX, x, y, leftW, boxH, "Used Filters")
	leftBox.SetAlign(fltk.ALIGN_CENTER | fltk.ALIGN_INSIDE)
	leftList := fltk.NewMultiBrowser(x, y+boxH, leftW, listH, "")

	// Right List: Available Items
	rightBox := fltk.NewBox(fltk.NO_BOX, x+w-rightW-10, y, rightW, boxH, "Available Filters")
	rightBox.SetAlign(fltk.ALIGN_CENTER | fltk.ALIGN_INSIDE)
	rightList := fltk.NewMultiBrowser(x+w-rightW-10, y+boxH, rightW, listH, "")

	// Arrow Buttons
	moveLeftButton := fltk.NewButton(btnX, btnY-40, btnW, btnH, "<--")  // Move from right to left
	moveRightButton := fltk.NewButton(btnX, btnY+10, btnW, btnH, "-->") // Move from left to right

	dlb := &DualListBox{
		Group:           group,
		leftBox:         leftBox,
		rightBox:        rightBox,
		leftList:        leftList,
		rightList:       rightList,
		moveLeftButton:  moveLeftButton,
		moveRightButton: moveRightButton,
		rightItems:      nil,
		leftItems:       nil,
		x:               x,
		y:               y,
		w:               w,
		h:               h,
	}

	// Move selected item from right list to left list
	moveLeftButton.SetCallback(func() {
		selected := rightList.Value()
		if selected > 0 {
			item := rightList.Text(selected)
			dlb.leftItems = append(dlb.leftItems, item)
			// Remove from rightItems
			index := slices.Index(dlb.rightItems, item)
			if index >= 0 {
				dlb.rightItems = slices.Delete(dlb.rightItems, index, index+1)
			}
			// Sort both lists and refresh UI
			sort.Strings(dlb.leftItems)
			sort.Strings(dlb.rightItems)
			dlb.Refresh()
			if dlb.onMoveLeft != nil {
				dlb.onMoveLeft()
			}
		}
	})

	// Move selected item from left list to right list
	moveRightButton.SetCallback(func() {
		selected := leftList.Value()
		if selected > 0 {
			item := leftList.Text(selected)
			dlb.rightItems = append(dlb.rightItems, item)
			index := slices.Index(dlb.leftItems, item)
			if index >= 0 {
				dlb.leftItems = slices.Delete(dlb.leftItems, index, index+1)
			}
			// Sort both lists and refresh UI
			sort.Strings(dlb.leftItems)
			sort.Strings(dlb.rightItems)
			dlb.Refresh()
			if dlb.onMoveRight != nil {
				dlb.onMoveRight()
			}
		}
	})

	group.End()
	dlb.Refresh()
	return dlb
}

// SetLeftTitle changes the label above the left list.
//
// title: The label to display.
func (d *DualListBox) SetLeftTitle(title string) {
	if d.leftBox != nil {
		d.leftBox.SetLabel(title)
	}
}

// SetRightTitle changes the label above the right list.
//
// title: The label to display.
func (d *DualListBox) SetRightTitle(title string) {
	if d.rightBox != nil {
		d.rightBox.SetLabel(title)
	}
}

// RegisterMoveRightHandler sets the callback to invoke after an item is moved from left to right.
//
// cb: The function to call after move right.
func (d *DualListBox) RegisterMoveRightHandler(cb func()) {
	d.onMoveRight = cb
}

// RegisterMoveLeftHandler sets the callback to invoke after an item is moved from right to left.
//
// cb: The function to call after move left.
func (d *DualListBox) RegisterMoveLeftHandler(cb func()) {
	d.onMoveLeft = cb
}

// Refresh clears and repopulates both lists from the current leftItems and rightItems slices.
// Both lists are displayed in sorted order.
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

// Resize moves and resizes the DualListBox and all its child widgets.
// This method should be called when the parent is resized.
func (d *DualListBox) Resize(x, y, w, h int) {
	d.x, d.y, d.w, d.h = x, y, w, h
	d.Group.Resize(x, y, w, h)

	// Re-layout children (fixed widths for lists/buttons, adjust positions)
	btnW := 50
	btnH := 30
	btnX := x + w/2 - btnW/2
	btnY := h / 2 // vertical center for buttons
	leftW := (w+btnW)/2 - 20
	rightW := (w+btnW)/2 - 20
	listH := h - 30
	boxH := 30

	// Left list and label
	d.leftList.Resize(x, y+boxH, leftW, listH)
	if d.leftBox != nil {
		d.leftBox.Resize(x, y, leftW, boxH)
	}

	// Right list and label
	rightX := x + w - rightW
	d.rightList.Resize(rightX, y+boxH, rightW, listH)
	if d.rightBox != nil {
		d.rightBox.Resize(rightX, y, rightW, boxH)
	}

	// Buttons between lists
	d.moveLeftButton.Resize(btnX, btnY-40, btnW, btnH)
	d.moveRightButton.Resize(btnX, btnY+10, btnW, btnH)
}

// SetLeftItems replaces the items in the left list (used/selected items).
// The list will be displayed in sorted order.
func (d *DualListBox) SetLeftItems(items []string) {
	d.leftList.Clear()
	sort.Strings(items)
	for _, item := range items {
		d.leftList.Add(item)
	}
	d.leftItems = slices.Clone(items)
}

// GetLeftItems returns all items currently in the left list (used/selected items).
func (d *DualListBox) GetLeftItems() []string {
	var items []string
	for i := 1; i <= d.leftList.Size(); i++ {
		items = append(items, d.leftList.Text(i))
	}
	return items
}

// SetRightItems replaces the items in the right list (available items).
// The list will be displayed in sorted order.
func (d *DualListBox) SetRightItems(items []string) {
	d.rightList.Clear()
	sort.Strings(items)
	for _, item := range items {
		d.rightList.Add(item)
	}
	d.rightItems = slices.Clone(items)
}

// GetRightItems returns all items currently in the right list (available items).
func (d *DualListBox) GetRightItems() []string {
	var items []string
	for i := 1; i <= d.rightList.Size(); i++ {
		items = append(items, d.rightList.Text(i))
	}
	return items
}

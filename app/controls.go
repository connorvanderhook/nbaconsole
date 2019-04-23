package app

import (
	"github.com/jroimartin/gocui"
)

func keybindings(g *gocui.Gui) error {
	var err error
	if err = g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, MoveUp); err != nil {
		return err
	}

	if err = g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, MoveDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, nil); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding(scoreboardLabel, gocui.KeyCtrlO, gocui.ModNone, nba.getBoxScore); err != nil {
		return err
	}

	return nil
}

// MoveDown modifies the gocui cursor +1 on y axis
func MoveDown(g *gocui.Gui, v *gocui.View) error {
	_, y := g.CurrentView().Cursor()
	v.SetCursor(0, y+1)
	return nil
}

// MoveUp modifies the gocui cursor +1 on y axis
func MoveUp(g *gocui.Gui, v *gocui.View) error {
	_, y := g.CurrentView().Cursor()
	v.SetCursor(0, y-1)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

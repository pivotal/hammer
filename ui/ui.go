package ui

import (
	"fmt"
	"io"
)

type UI struct {
	Out io.Writer
	Err io.Writer
}

func (ui *UI) DisplayText(text string) {
	fmt.Fprint(ui.Out, text)
}

func (ui *UI) DisplayError(err error) {
	fmt.Fprint(ui.Err, err.Error())
}

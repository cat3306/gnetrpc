package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/cat3306/gnetrpc/tool/debugtool/the"
	"github.com/cat3306/gnetrpc/tool/debugtool/view"
)

var (
	width  float32 = 800
	height float32 = 600
)

func main() {
	a := app.New()
	a.Settings().SetTheme(&the.MyTheme{})
	w := a.NewWindow("gnetrpc debug tool")
	w.Resize(fyne.NewSize(width, height))
	upperLeft := view.InitUpperLeftView()
	view.GlobalText = view.NewTextView()
	view.GlobalText.Daemon()
	left := container.NewVSplit(container.NewVBox(
		upperLeft.Join(),
	), view.GlobalText.Join())

	view.ReceiveText = view.NewTextView()
	view.ReceiveText.Daemon()
	right := container.NewVSplit(container.NewVBox(
		new(view.SendView).Init(),
	), view.ReceiveText.Join())
	all := container.NewHSplit(left, right)

	w.SetContent(all)

	w.ShowAndRun()
}

package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func fetchAndDisplayUS30YData(content *fyne.Container, FetchDataUS30Y func() ([][]string, error), placeholder fyne.CanvasObject) {
	go func() {
		us30yData, err := FetchDataUS30Y()
		if err != nil {
			content.Objects = replaceObject(content.Objects, placeholder, widget.NewLabel("Error fetching US30Y data"))
		} else {
			us30yCard := renderUS30YCard(us30yData)
			content.Objects = replaceObject(content.Objects, placeholder, us30yCard)
		}
		content.Refresh()
	}()
}

func renderUS30YCard(data [][]string) *widget.Card {
	return widget.NewCard(
		"US30Y Data",
		"",
		container.NewVBox(
			widget.NewLabel("Name: " + data[1][0]),
			widget.NewLabel("Last: " + data[1][1]),
			widget.NewLabel("Open: " + data[1][2]),
			widget.NewLabel("High: " + data[1][3]),
			widget.NewLabel("Low: " + data[1][4]),
		),
	)
}

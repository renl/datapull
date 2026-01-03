package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func fetchAndDisplaySP500Data(content *fyne.Container, FetchDataSP500 func() ([]string, error), placeholder fyne.CanvasObject) {
	go func() {
		sp500Data, err := FetchDataSP500()
		if err != nil {
			content.Objects = replaceObject(content.Objects, placeholder, widget.NewLabel("Error fetching S&P 500 data"))
		} else {
			sp500Card := renderSP500Card(sp500Data)
			content.Objects = replaceObject(content.Objects, placeholder, sp500Card)
		}
		content.Refresh()
	}()
}

func renderSP500Card(data []string) *widget.Card {
	if len(data) < 5 {
		return widget.NewCard(
			"S&P 500 Index",
			"",
			widget.NewLabel("Insufficient data"),
		)
	}

	currentIndex := data[0]
	previousClose := data[1]
	volume := data[2]
	lastWeekHigh := data[3]
	lastWeekLow := data[4]

	return widget.NewCard(
		"S&P 500 Index",
		"",
		container.NewVBox(
			widget.NewLabel(currentIndex),
			widget.NewLabel(previousClose),
			widget.NewLabel(volume),
			widget.NewLabel(lastWeekHigh),
			widget.NewLabel(lastWeekLow),
		),
	)
}

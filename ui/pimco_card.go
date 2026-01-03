package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func fetchAndDisplayPimcoData(content *fyne.Container, FetchDataPimcoGISIncome func() ([]string, error), placeholder fyne.CanvasObject) {
	go func() {
		pimcoData, err := FetchDataPimcoGISIncome()
		if err != nil {
			content.Objects = replaceObject(content.Objects, placeholder, widget.NewLabel("Error fetching Pimco data"))
		} else {
			pimcoCard := renderPimcoCard(pimcoData)
			content.Objects = replaceObject(content.Objects, placeholder, pimcoCard)
		}
		content.Refresh()
	}()
}

func renderPimcoCard(data []string) *widget.Card {
	if len(data) < 3 {
		return widget.NewCard(
			"Pimco GIS Income",
			"",
			widget.NewLabel("Insufficient data"),
		)
	}

	// companyName := data[0]
	currentPrice := data[1]
	lastTradingDayPrice := data[2]

	return widget.NewCard(
		"Pimco GIS Income",
		"",
		container.NewVBox(
			// widget.NewLabel("Company Name: " + companyName),
			widget.NewLabel("Current Price: " + currentPrice),
			widget.NewLabel("Last Trading Day Price: " + lastTradingDayPrice),
		),
	)
}
